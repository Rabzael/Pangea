# Pangea
> A simple and simplified forward-proxy-caching-IDS-AI-log-analyzer comprehensive system, made for self-educational purpose.

Long story short: with the help of ChatGPT I found some ideas to practice my *coding skills* (!) and learn something new about networks, data, AI...

LLMs (mostly [ChatGPT](https://chat.openai.com/) and [Claude](https://claude.ai/)) have been used to analyze and design the, generate/study parts of the code and to write part of the documentation.

## Modules
The whole project is divided in modules. All the documentation in the ['docs'](./docs/) folder:

| #                         | Name     | Description                   | Languages |
|---------------------------|----------|-------------------------------|-----------|
| [00](docs/00-overview.md) | Overview | Architecture and instructions | -         |
| [01](docs/01-proxy.md)    | Proxy    | A forward proxy               | Go        |
| [02](docs/02-etl.md)      | ETL      | Log ETL tool                  | Rust      |

Every module has its own **Dockerfile** in the [docker](./docker/) folder. Container orchestration is done via **Kubernetes** using the files in the [k8s](./k8s/) folder. 
Communication between modules is done via files and **Kafka** topics.
More about general architecture and usage in the [Overview](docs/00-overview.md).

## Special thanks to
My beloved wife and my wonderful kids, who let me some spare time here and there to purse my hobbies.

## Resources
### Generic
- RFC 2616 on HTTP/1.1 Protocol, especially [section 9](https://www.rfc-editor.org/rfc/rfc2616#section-9)
- [Arrow format documentation](https://arrow.apache.org/docs/format/Columnar.html)

### Go
- [Go by Example](https://gobyexample.com/) contains all you need to learn/refresh Go
- Then the [http package docs](https://pkg.go.dev/net/http@go1.25.2) is mandatory
- Alex Rios, [_System Programming Essentials with Go_](https://www.packtpub.com/en-it/product/system-programming-essentials-with-go-9781801813440), Packt Publishing 2024

### Rust
- [The Rust Programming Language](https://doc.rust-lang.org/book/title-page.html)
- [docs.rs](https://docs.rs/) for all the crates documentation

## Roadmap
Subject to changes.

| Task                                      | Module  | Status      |
|-------------------------------------------|---------|-------------|
| Implement a forward proxy server          | Proxy   | Completed   |
| Implement a log ETL tool                  | ETL     | Completed   |
| Implement Kafka communication             | -       | In Progress |
| Containerize the modules with Docker      | -       | In Progress |
| Orchestrate containers with Kubernetes    | -       | In Progress |
| Implement AI-based log analysis           | Agent   | Pending     |
| Implement IDS functionalities             | IDS     | Pending     |
| Implement blocklist functionalities       | Proxy   | Pending     |
| Implement caching functionalities         | Proxy   | Pending     |

## License
This project is under Apache 2.0 License.