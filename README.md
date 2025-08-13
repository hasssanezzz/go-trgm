# Trigram Index

**Status:** 🚧 Work in Progress

Trigram Index is an experimental Go-based indexing engine designed for **fast substring search** using **trigram indexing**. It works by breaking strings into fixed-size 3-character chunks (trigrams), mapping them to integer IDs, and maintaining both **in-memory** and **on-disk** indexes for efficient lookups.

## Features (WIP)

* **In-memory trigram index** for fast search
* **Disk-backed storage** to handle large datasets beyond memory limits
* **Automatic flushing** of trigrams to disk when their occurrence count exceeds a threshold
* **Fetch API** that merges results from memory and disk
* Configurable threshold for flushing

## Current Status

* ✅ Basic indexing and fetching implemented
* ✅ Threshold-based disk flushing
* ✅ Search from both memory and disk
* ❌ No durability, blocks index is not yet serialized and deserialized
* ❌ No crash recovery
* ❌ No concurrency support yet
* ❌ No CLI / API layer yet
* ❌ No real-world performance benchmarks

## Getting Started

```bash
go get github.com/hasssanezzz/trigram-index
```

Example usage will be added soon as the project evolves.
