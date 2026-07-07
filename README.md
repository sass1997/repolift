# Repolift 🏋️‍♂️

> A declarative, lightning-fast workspace manager for your Git repositories.

[![Go Report Card](https://goreportcard.com/badge/github.com/sass1997/repolift)](https://goreportcard.com/report/github.com/sass1997/repolift)

Repolift allows you to define your local developer folder structure and git repositories declaratively via a simple YAML file. Instead of manually creating directories and cloning dozens of repositories across multiple projects, you just define your desired state, run one command, and Repolift does the heavy lifting.

## 🚀 Features (v1)

- **Declarative Configuration**: Define your workspaces and repositories in a clean `yaml` format.
- **Idempotent**: Run it as many times as you want. Repolift detects what's already cloned and only fetches what's missing.
- **Native Git Integration**: Uses your system's git, ensuring all your existing SSH keys, GPG signing setups, and git-configs work out of the box.
- **Fast & Lightweight**: Written in Go. Distributed as a single, dependency-free binary.

## 📦 Installation

*(Coming soon: Nix package & Homebrew tap)*

For now, you can build it from source:

```bash
git clone https://github.com/sass1997/repolift.git
cd repolift
go build -o repolift main.go
sudo mv repolift /usr/local/bin/
```

## 🛠️ Usage

Create a `repolift.yaml` file to define your workspace structure:

```yaml
workspaces:
  # The path where the repositories should be cloned. 
  # Tilde (~) expansion is fully supported.
  - path: "~/work/project-alpha"
    repositories:
      - url: "git@github.com:myorg/api-service.git"
        dir: "api"
      - url: "git@github.com:myorg/frontend-app.git"
        dir: "frontend"
        
  - path: "~/personal/open-source"
    repositories:
      - url: "https://github.com/spf13/cobra.git"
        dir: "cobra"
```

Apply the configuration:

```bash
repolift apply -f repolift.yaml
```

## 🗺️ Roadmap

- [x] **v1: Core Foundation** - Declarative folder structures and repository cloning.
- [ ] **v2: Advanced Authentication** - Per-repository SSH keys, credential encryption, and Git config management.
- [ ] **v3: Automation & Scheduling** - Background jobs to automatically pull the latest `main` branches (e.g., every morning at 7:00 AM) so you start your day with fresh code.

## 🤝 Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

## 📄 License

MIT License
