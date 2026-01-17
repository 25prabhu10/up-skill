# Up Skill

Learning to solve problems, data structures and algorithms.

## JavaScript/TypeScript

Requirements any JavaScript/TypeScript runtime is fine, but we recommend using:

- [Bun](https://bun.com)

Install dependencies:

```bash
bun install
```

Run the unit tests:

```bash
bun run test
```

Generate JavaScript:

```bash
bun run build
```

## C/C++

Requirements:

- [GCC](https://gcc.gnu.org/) (version 15 or higher)
- [Make](https://www.gnu.org/software/make/) (version 4 or higher)

Run the unit tests:

```bash
make test
```

Build files:

```bash
make build
```

Get help on available commands:

```bash
make help
```

## Python

Requirements:

- [uv](https://docs.astral.sh/uv/) (version 0.8 or higher)
- [Python](https://www.python.org/) (version 3.13 or higher)

Install dependencies:

```bash
uv python install

uv sync --all-extras
```

Run the unit tests:

```bash
uv run pytest
```
