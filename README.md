# DOTS-go

`DOTS-go` is a high-performance OSINT tool written in Go, designed for scraping
GitHub repositories to extract valuable metadata like contributor aliases, email
addresses, and commit history. It supports multithreading and various export
formats for easy integration into your OSINT workflow.

## Features

- **Multithreaded Execution**: Leverages Go goroutines to scale scraping tasks across multiple cores.
- **Multiple Input Methods**:
  - Scrape repositories directly from a list of URLs.
  - Process a text file containing repository links.
  - Scan a local directory for existing git repositories.
- **Versatile Export Options**: Support for STDOUT, TXT, CSV, XLS (Excel), XML, and JSON.
- **Data Cleaning**: Automatically filters out anonymous `users.noreply.github.com` email addresses.
- **⚡ Efficient Caching**: Manages a local cache for repository cloning and processing.

## Installation

Ensure you have [Go](https://go.dev/doc/install) (1.23.3+) installed.

```bash
git clone https://github.com/MisterNorwood/DOTS-go.git
cd DOTS-go
go build -o dots-go main.go
```

## Usage

Run the tool using the compiled binary:

```bash
./dots-go [flags]
```

### Flags

| Flag | Alias | Default | Description |
| :--- | :--- | :--- | :--- |
| `--threads` | `-t` | `4` | Number of threads for multithreaded workloads. |
| `--file` | `-f` | - | Plain text file containing repository links (one per line). |
| `--links` | `-l` | - | Direct links of repositories to be scraped (multiple allowed). |
| `--repoDir` | `-r` | - | Local directory containing git repositories to scrape. |
| `--exportForm`| `-e` | `TXT` | Export formats: `STDOUT`, `CSV`, `XLS`, `TXT`, `JSON`, `XML`, or `ALL`. |
| `--stripNoreply`| `-n` | `true` | Strip default anonymous github mails. |
| `--version` | `-v` | - | Print the version. |

### Examples

**Scrape specific links and export to JSON:**
```bash
./dots-go --links https://github.com/example/repo1,https://github.com/example/repo2 --exportForm JSON
```

**Process a list of repos from a file using 8 threads:**
```bash
./dots-go -f repos_list.txt -t 8 -e CSV
```

**Scan a local directory of repositories:**
```bash
./dots-go -r /home/user/my-repos -e XLS
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
