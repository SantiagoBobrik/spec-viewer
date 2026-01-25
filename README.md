# Spec Viewer
<p align="center">
  <img width="470" height="224" alt="image" src="https://github.com/user-attachments/assets/cff9bdc6-8461-4dcb-af3e-ef7cf1403873" />
  <img width="1439" height="754" alt="image" src="https://github.com/user-attachments/assets/1ba84741-d80d-42c3-ae49-a873ba35b8be" />
</p>

> [!IMPORTANT]
> This project is currently under active development. Features and interfaces are subject to change.

Spec Viewer is a specialized visualization tool designed to accompany the **Spec Kit** ecosystem. It facilitates **Spec Driven Development (SDD)** by providing a high-fidelity, live preview of technical specifications as they evolve.

Unlike generic Markdown viewers, Spec Viewer is tailored to the specific needs of the SDD workflow, ensuring that specifications remain the readable and authoritative source of truth for your development process. It serves as the visual interface for your "living specifications," allowing developers and stakeholders to review architectural decisions and data models in real-time.



## Features

- **SDD Optimization**: Designed to render Spec Kit artifacts with precision.
- **Live Synchronization**: Instant feedback loop for file changes using efficient WebSocket connections.
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
