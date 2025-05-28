# Declarative Configurator

**Declarative Configurator** is a modular system configuration tool written in Go, designed to bring declarative infrastructure principles to personal system management. Through simple YAML files, users can define desired state for things like installed packages and UI settings (coming soon), and apply them with a single command.


> ⚠️ **Linux-only:** This tool is currently built for Linux systems.

---

## ✨ Features

- ⚙️ Declarative configuration using easy-to-read YAML
- 📦 Package management module (currently supports `dnf`)
- 🧩 Modular architecture for easy extension (APT, Flatpak, UI configs coming soon)
- 🔁 Simple CLI with intuitive commands

---

## 🚀 Usage

```bash
# Refresh all modules
declarative-configurator refresh all

# Refresh only the packages module
declarative-configurator refresh packages
```

---

## 📝 Configuration Example

Place a YAML file in the expected config path (or pass via flag, if supported):

```yaml
fedora:
  Native:
    - btop
  Flatpaks:
  Local:
```

---

## 🧠 Available Commands

| Command            | Description                             |
|--------------------|-----------------------------------------|
| `refresh`          | Alias for `refresh all`                 |
| `refresh all`      | Refresh all configured modules          |
| `refresh packages` | Refresh only the packages module        |
| `update packages`  | Update All Packages (Native, Flatpak)   |

---

## 🗺 Roadmap

* [x] DNF support
* [x] APT support
* [x] Flatpak support
* [x] Local Packages support
* [x] Package Updates (All)
* [ ] Package Update (Individual)
* [ ] UI customization via YAML
* [ ] Per-module config validation and schema checking

---

## 🤝 Contributing

Contributions welcome! Open issues, submit PRs, or suggest enhancements.

If you'd like to add a new module, the modular design allows for easy expansion. Documentation for creating new modules is coming soon.

---

## 🙌 Acknowledgments

Inspired by tools like **Ansible** and **Nix** — giants in the world of system configuration. This project doesn't aim to compete with them, but rather to serve a simpler, more personal purpose: making it easier for me to spin up fresh Fedora instances with minimal hassle. It's simple, and it does what I need for now — but I do hope to gradually improve it and maybe add more features along the way.
