# Spec Viewer
<p align="center">
<img width="1622" height="458" alt="image" src="https://github.com/user-attachments/assets/6aec476d-6ea9-46f3-b45e-6dbbe8189cf0" />
<img width="1874" height="312" alt="image" src="https://github.com/user-attachments/assets/5fc1c851-b6dd-4ae5-adb6-6b9587469865" />
<img width="1433" height="689" alt="image" src="https://github.com/user-attachments/assets/6f61f055-273b-4435-a975-d9b3d8e0003f" />
<img width="1433" height="1433" alt="image" src="https://github.com/user-attachments/assets/87db29d2-6e30-4e2b-8f80-7c8f4eeeddea" />



</p>

> [!IMPORTANT]
> This project is currently under active development. Features and interfaces are subject to change.

Spec Viewer is a specialized visualization tool designed to accompany the **Spec Kit** ecosystem. It facilitates **Spec Driven Development (SDD)** by providing a high-fidelity, live preview of technical specifications as they evolve.

Unlike generic Markdown viewers, Spec Viewer is tailored to the specific needs of the SDD workflow, ensuring that specifications remain the readable and authoritative source of truth for your development process. It serves as the visual interface for your "living specifications," allowing developers and stakeholders to review architectural decisions and data models in real-time.



## Features

- **SDD Optimization**: Designed to render Spec Kit artifacts with precision.
- **Live Synchronization**: Instant feedback loop for file changes using WebSocket connections with scroll-preserving hot reload.
- **GitHub Flavored Markdown**: Full support for tables, task lists, strikethrough, and auto-linked URLs.
- **Mermaid Diagrams**: Render flowcharts, sequence diagrams, ER diagrams, and more directly in your specs.
- **Table of Contents**: Auto-generated from headings with desktop sidebar and mobile overlay.
- **Sidebar Search**: Filter specs by file or folder name.
- **Mobile Responsive**: Collapsible sidebar and TOC overlays for mobile and tablet.
- **Zero Configuration**: Adheres to Spec Kit conventions "out of the box" without requiring complex setup.
- **Global Accessibility**: Runs as a standalone CLI tool primarily for local development environments.

## Installation

Spec Viewer is distributed as a Go binary and should be installed globally to be accessible from any project directory.

Ensure you have Go installed (version 1.24 or higher).

```bash
go install github.com/SantiagoBobrik/spec-viewer/cmd/spec-viewer@latest
```

## Usage

To start the viewer, navigate to the root of your project or the directory containing your specifications:

```bash
spec-viewer serve
```

By default, the server listens on port `9091` and watches the `./specs` directory, following standard Spec Kit structure.

### Configuration

You can override the default behaviors using command-line flags:

| Flag | Shorthand | Description | Default |
|------|-----------|-------------|---------|
| `--port` | `-p` | Port to run the server on | `9091` |
| `--folder` | `-f` | Directory to watch for Markdown files | `./specs` |

### Workflow Example

1. Generate specifications using Spec Kit.
2. Run `spec-viewer serve` in a separate terminal window.
3. Open `http://localhost:9091` in your browser.
4. As you or your AI agents update the specifications, the viewer will automatically refresh to reflect the latest state.

## Contributing

This project is open source and welcomes contributions. Please ensure all pull requests adhere to the existing architectural standards.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
