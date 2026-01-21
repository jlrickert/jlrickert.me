---
date: "2025-12-04T22:52:01-06:00"
draft: false
title: "What Is a Keg"
slug: what-is-a-keg
description: "Understanding Knowledge Exchange Graphs (KEGs)—personal knowledge management systems that organize information as interconnected nodes instead of traditional folders."
---

# What is a KEG? Building Your Own Personal Knowledge Graph

You've probably heard the term "knowledge base" before. It usually refers to a collection of documents, articles, or FAQ
pages. But what if your knowledge could be organized more like your brain actually works—as a network of interconnected
ideas rather than a linear collection of files?

That's what a KEG is. It stands for **Knowledge Exchange Graph**, and it's a personal knowledge management system that
organizes information as a network of interconnected nodes instead of traditional folders and documents.

## From Files to Graphs

Traditional knowledge management goes like this:

- Create a document
- File it in a folder
- Search for it by name when you need it

A KEG works differently:

- Create a node (a discrete piece of knowledge)
- Tag it with relevant categories (domains)
- Link it to related nodes
- Discover it through multiple paths: direct lookup, tag searches, or by following connections

The difference is profound. With files, your knowledge is isolated. With a KEG, your knowledge is **connected**.

## The Building Blocks

A KEG is built on two core concepts:

### Ontology: What Kind of Thing Is This?

An ontology defines the types of entities that can exist in your knowledge base. It answers: "What _type_ of thing is
this?"

Common entity types might include:

- **Projects** - Long-running initiatives
- **Issues** - Problems to solve
- **Concepts** - Abstract ideas
- **Hardware** - Physical equipment
- **Recipes** - Cooking instructions
- **People** - Notable individuals
- **Patches** - Modifications to things
- **Tricks** - Technical tips and techniques

Here's the key insight: **Each KEG defines its own ontology.** There's no universal standard. Your personal KEG might
emphasize hardware and infrastructure, while someone else's focuses on creative writing and storytelling. A KEG makes no
specification for the hard existence of any particular entity types—you design the ontology that makes sense for your
way of thinking.

### Domains: What Field Does This Belong To?

Domains are conceptual categories that organize knowledge by subject area or field. They answer: "What _area of
knowledge_ does this belong to?"

If an ontology says "this is a Project," a domain says "this project is related to homelab infrastructure, DevOps,
networking, and DNS."

A single node can belong to multiple domains. This creates natural clustering in your knowledge base—all "homelab"
related content clusters together, making it easy to discover related concepts.

## How They Work Together

Node: "Current Homelab Setup"\
├── Ontology: Project\
└── Domains: homelab, sysadmin, devops, networking, virtualization

This node is a **Project** (ontology) that spans multiple knowledge areas: homelab infrastructure, system
administration, DevOps practices, networking, and virtualization (domains).

When you later query "show me everything tagged with 'networking'," this node appears—along with DNS guides, router
configurations, VPN setup notes, and anything else related to networking.

## Why a Graph Matters

Here's what makes a KEG different from a traditional note-taking app:

| Aspect            | Traditional Notes           | Knowledge Graph                         |
| ----------------- | --------------------------- | --------------------------------------- |
| **Organization**  | Hierarchical folders        | Network of connections                  |
| **Discovery**     | Search-dependent            | Multiple paths (tags, links, backlinks) |
| **Relationships** | Hidden connections          | Explicit and visible                    |
| **Growth**        | More files = harder to find | More nodes = more powerful              |
| **Serendipity**   | Rare                        | Frequent through exploration            |

With a KEG, as you add more knowledge, the system becomes more useful. Unexpected connections surface. You discover
relationships you didn't know existed.

## The Power of Multiple Query Paths

A KEG supports several ways to find knowledge:

**Direct lookup** - If you remember the concept ID, access it instantly

**Tag-based discovery** - "Show me all nodes about homelab" - returns everything tagged with that domain

**Associative search** - Start with one node and follow links to related concepts

**Content search** - Search across all node content when you can't remember the category

**Reverse traversal** - See what other nodes reference a particular concept (backlinks)

This flexibility means you don't have to remember the perfect folder structure or exact filename. Multiple retrieval
paths mean you can rediscover your knowledge in different ways.

## A Simple Implementation

A KEG doesn't require fancy software. The core is remarkably simple:

- **Numbered nodes** - Each piece of knowledge gets a unique ID (1, 2, 3, etc.)
- **Markdown files** - Content is just markdown with a README.md file
- **Tags** - Metadata in YAML files for categorization
- **Links** - References between nodes using simple markdown links
- **Command-line queries** - Simple tools to search, filter, and discover

No databases. No proprietary formats. Just files, tags, and connections.

## Who Should Build a KEG?

KEGs work well for:

- **Software engineers** wanting to document projects, architecture, and technical knowledge
- **System administrators** managing infrastructure and operations knowledge
- **Researchers** organizing papers, findings, and related concepts
- **Makers and hobbyists** tracking projects, techniques, and ideas
- **Anyone** who wants to own their knowledge and make it interconnected

The advantage is that you control the structure. Your KEG reflects how _you_ think, not how a SaaS company thinks you
should organize information.

## Beyond Personal: Knowledge Exchange

The "Exchange" in KEG hints at something deeper. While a KEG starts personal, the structure enables sharing. Your nodes
can link to other people's KEGs. Your domains can align with others' domains. Knowledge can be shared while remaining
under your control.

## Getting Started

If you want to build your own KEG:

1. **Define your ontology** - What types of things will you store? (Projects, Issues, Articles, Concepts, etc.)
2. **Identify your domains** - What subject areas matter to you? (programming, infrastructure, cooking, etc.)
3. **Create nodes** - Start capturing knowledge as discrete, atomic pieces
4. **Tag consistently** - Apply your domains to create clustering
5. **Link thoughtfully** - Connect related nodes to build the graph
6. **Query and discover** - Use simple tools to explore your knowledge

You don't need to be perfect from day one. Your ontology and domains will evolve as you use them.

## The Future of Personal Knowledge

We're moving beyond the era of note-taking apps that emphasize capture. The future is about knowledge systems that
emphasize **discovery and emergence**. Systems where the more you add, the smarter it becomes. Systems that make
unexpected connections visible.

A KEG is one approach to that future—a personal knowledge management system that's simple enough to be portable,
flexible enough to reflect your thinking, and powerful enough to create genuine serendipitous discovery.

Your knowledge deserves to be more than a folder structure. It deserves to be a graph.
