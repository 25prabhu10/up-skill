# Up Skill

Learning to solve problems, data structures and algorithms.

## JavaScript/TypeScript

Requirements:

- [Node.js](https://nodejs.org/) (version 22 or higher)
- [pnpm](https://pnpm.io/) (version 10 or higher)

Install dependencies:

```bash
pnpm install
```

Run the unit tests:

```bash
pnpm run test
```

Generate JavaScript:

```bash
pnpm run build
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
