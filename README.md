# Lifegame Multiverse
**Conway's Game of Life** is something that has been close to me since the start of my degree. 

In my first year I implemented a solution in Java for a programming project (which had an UI but we only did the logic) and I kind of became the meme because of the amount of nested for loops I used (4) which I think was a nice try.

Then I wanted to prove them wrong and did an implementation in C with only one for loop and showed the game in the terminal, which, by looking at it some year later, I noticed it leaked a lot of memory.

Now I come back to it, inspired by [carykh's video](https://www.youtube.com/watch?v=QK_KZv-YyOc), with an expanded perception of the game and with a new language: Go.

Wellcome then, to my third (and I hope final) attempt at this beautiful game.

This time made to work with every possible set of rules and with a graphical UI thanks to [Ebitengine](https://ebitengine.org/).

## Purpose
The purpose of this project is to improve my knowledge of Go and of how to structure a Go module properly.

After reading [Organizing a Go module](https://go.dev/doc/modules/layout)
I decided to use the **Multiple commands** approach, even though I have only one `main.go` but I like to have it under the `lifegame/` folder.

The actual structure is:
```
lifegame-multiverse/
├── go.mod
├── go.sum
├── internal
│   ├── rules
│   └── world
├── LICENSE
├── lifegame
│   └── main.go
└── README.md
```
Might be too complicated for the size of the project but I want to experiment well with handling multiple packages.

## Requirements
This attempt will use `goroutines` to parallelize the next generations, so the plan is to be able to handle **BIG** universes or even multiple medium-size universes at a time.

So the requirement list will be:
- [ ] Can handle really big universes at 60fps.
- [ ] Camera placement and zoom control in the UI to see and navigate big universes.
- [ ] Can simulate various universes at the same time in the same space.

Let's see what I am able to achieve.

### Update 31/8/2026
The program can simulate a single 2560x1440 world fullscreen in my laptop (2025x1139) at 60TPS and 60FPS.\
That is checking, updating and showing 3.686.400 cells 60 times per second.\
Is this a really big universe? For my current hardware I think it is, but I will try to squeeze a little bit before marking the first checkbox.
