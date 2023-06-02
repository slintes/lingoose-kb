# Using OpenAI for querying a single Wikipedia page

## Usage

1. export your OpenAI API key:

```shell
export OPENAI_API_KEY="your_key"
```

2. Get the title of the Wikipedia page you want to search:

e.g. for `https://en.wikipedia.org/wiki/Red_Hat` the title is `Red_Hat`

3. Run:

```shell
$ go run github.com/slintes/lingoose-kb "<title>"
```

4. Chat (example for `"Fencing_(computing)"`)

```shell
Learning Wiki page... Done
Enter a query to search the knowledge base. Type 'quit' to exit.
> What is node isolation?
Node isolation is the process of ensuring that I/O can no longer be done from a malfunctioning node in a computer cluster in order to protect processes from other active nodes modifying the resources during node failures. This is typically done automatically by cluster infrastructure such as shared disk file systems to prevent a split brain condition and data corruption.
> What happens when the number of nodes increases?
When the number of nodes in a cluster increases, the likelihood that one of them may fail at some point also increases.
> quit
```

5. When using a new title, delete the `./db.json` file first!
