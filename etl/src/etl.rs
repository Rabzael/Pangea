use std::path::Path;
use std::fs::File;
use std::io::{BufRead, BufReader};
use std::sync::Arc;
use serde_json::{Value, Map};
use arrow::ipc::writer::FileWriter;
use arrow::datatypes::{Schema, Field, DataType};
use arrow::record_batch::RecordBatch;
use arrow::array::StringArray;

fn main() {  
    // Check arguments
    let args: Vec<String> = std::env::args().collect();
    if args.len() < 4 {
        println!("Usage: {} <batch size> <input_logfile> <output_file>", args[0]);
        std::process::exit(1);
    }
    let batch_size: usize = match args[1].parse() {
        Ok(size) => size,
        Err(_) => {
            println!("Invalid batch size");
            std::process::exit(1);
        }
    };
    let input_file = &args[2];
    let output_path = &args[3];

    // Check if input file and output directory exist
    if !Path::new(input_file).exists() {
        println!("Input file does not exist");
        std::process::exit(1);
    }
    if !Path::new(output_path).parent().unwrap_or(Path::new(".")).exists() {
        println!("Output directory does not exist");
        std::process::exit(1);
    }

    println!("Starting ETL process from \"{}\" to \"{}\"...", input_file, output_path);

    // Define readers
    let file = File::open(input_file).expect("Could not open input file");
    let reader = BufReader::new(file);

    // Define Arrow schema
    let schema = Arc::new(Schema::new(vec![
        Field::new("time", DataType::Utf8, false),
        Field::new("remote_addr", DataType::Utf8, false),
        Field::new("method", DataType::Utf8, false),
        Field::new("url", DataType::Utf8, false),
        Field::new("status", DataType::Utf8, false),
        Field::new("error", DataType::Utf8, true),
        Field::new("blocked", DataType::Utf8, true),
    ]));

    // Lines
    let mut time_lines = Vec::new();
    let mut remote_addr_lines = Vec::new();
    let mut method_lines = Vec::new();
    let mut url_lines = Vec::new();
    let mut status_lines = Vec::new();
    let mut error_lines = Vec::new();
    let mut blocked_lines = Vec::new();

    let mut line_count = 0;

    // Output Arrow file
    let output_file = File::create(output_path).expect("Could not create output file");
    let mut writer = match FileWriter::try_new(&output_file, &schema) {
        Ok(w) => w,
        Err(e) => {
            eprintln!("Could not create Arrow writer: {}", e);
            std::process::exit(1);
        }
    };

    // Read input line by line
    for line in reader.lines() {
        let buf: String = match line {
            Ok(l) => l.trim().to_string(),
            Err(e) => {
                eprintln!("Could not read line {}: {}", line_count, e);
                continue;
            }
        };

        line_count += 1;

        // Process data: convert to JSON object
        let mut json_data: Value = match serde_json::from_str(&buf) {
            Ok(data) => data,
            Err(e) => {
                eprintln!("{}: Could not parse JSON: {}", line_count, e);
                continue;
            }
        };

        // Process data: get mutable object
        let cleaned_data: &mut Map<String, Value> = match json_data.as_object_mut() {
            Some(obj) => obj,
            None => {
                eprintln!("{}: Invalid JSON object", line_count);
                continue;
            }
        };

        // Process data: clean unnecessary fields
        cleaned_data.remove("content_len");
        
        // Extract data
        push_value(&mut time_lines, cleaned_data, "time");
        push_value(&mut remote_addr_lines, cleaned_data, "remote_addr");
        push_value(&mut method_lines, cleaned_data, "method");
        push_value(&mut url_lines, cleaned_data, "url");
        push_value(&mut status_lines, cleaned_data, "status");
        push_value(&mut error_lines, cleaned_data, "error");
        push_value(&mut blocked_lines, cleaned_data, "blocked");

        // Write if limit reached
        if time_lines.len() >= batch_size {
            write_arr(
                &schema,
                &mut time_lines,
                &mut remote_addr_lines,
                &mut method_lines,
                &mut url_lines,
                &mut status_lines,
                &mut error_lines,
                &mut blocked_lines,
                &mut writer
            );
        }
    }

    // Write remaining data
    if !time_lines.is_empty() {
        write_arr(
            &schema,
            &mut time_lines,
            &mut remote_addr_lines,
            &mut method_lines,
            &mut url_lines,
            &mut status_lines,
            &mut error_lines,
            &mut blocked_lines,
            &mut writer
        );
    };

    // End
    writer.finish().expect("Could not finish writing");
    println!("ETL process completed. Processed {} lines.", line_count);
}


// Helper function to push value into array
fn push_value(array: &mut Vec<Option<String>>, data: &Map<String, Value>, key: &str) {
    let val: Option<String> = match data.get(key) {
        Some(v) => {
            match v {
                Value::String(s) => Some(s.clone()),
                _ => None
            }
        },
        None => None
    };
    array.push(val);
}

// Helper function to write Arrow arrays
fn write_arr(
    schema: &Arc<Schema>, 
    time_lines: &mut Vec<Option<String>>, 
    remote_addr_lines: &mut Vec<Option<String>>, 
    method_lines: &mut Vec<Option<String>>, 
    url_lines: &mut Vec<Option<String>>, 
    status_lines: &mut Vec<Option<String>>, 
    error_lines: &mut Vec<Option<String>>, 
    blocked_lines: &mut Vec<Option<String>>, 
    writer: &mut FileWriter<&File>,
) {
    let arr_writer = RecordBatch::try_new(
        schema.clone(),
        vec![
            Arc::new(StringArray::from(time_lines.clone())),
            Arc::new(StringArray::from(remote_addr_lines.clone())),
            Arc::new(StringArray::from(method_lines.clone())),
            Arc::new(StringArray::from(url_lines.clone())),
            Arc::new(StringArray::from(status_lines.clone())),
            Arc::new(StringArray::from(error_lines.clone())),
            Arc::new(StringArray::from(blocked_lines.clone())),
        ],
    ).expect("Could not create RecordBatch");

    writer.write(&arr_writer).expect("Could not write RecordBatch");

    time_lines.clear();
    remote_addr_lines.clear();
    method_lines.clear();
    url_lines.clear();
    status_lines.clear();
    error_lines.clear();
    blocked_lines.clear();
}