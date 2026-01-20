# lolzy

**lolzy** is a command-line tool (CLI) that helps you explore the current *League of Legends* patch meta. It uses a **TUI** (Terminal User Interface) and allows you to quickly check:

* which champions are strong picks in a given **role**,
* which champions are good **counters** against a specific champion.

## Features

lolzy provides two main commands:

* `lolzy meta <role> --options`
* `lolzy counter <champ> --options`

---

## `lolzy meta <role>`

This command shows, for a given **role** (e.g. top, jungle, mid, adc, support):

* which champions are considered strong in the current patch,
* and which champions perform well **against them** (counter information).

### Available options (`meta options`)

* `--all`, `-a`
  If enabled, champions with low sample size are also included (e.g. fewer than 10 observed matches or pick rate below 0.5%).

* `--top <number>`, `-t <number>`
  Limits the output to the **top X champions** only.

* `--role <role>`, `-r <role>`
  Specifies in which role counters should be searched (e.g. `top`, `jungle`).

* `--rank <rank>`, `-rk <rank>`
  Filters the data to matches played in the given rank.
  Examples: `bronze`, `silver`, `gold`, etc.

### Example

```bash
lolzy meta top -t 5 -rk silver
```

This command displays the **top 5 champions for top lane** based on **silver-ranked** matches.

---

## `lolzy counter <champ>`

This command shows which champions perform well **against a specific champion**.

### Available options (`counter options`)

* `--rank <rank>`, `-rk <rank>`
  Uses only matches from the specified rank.

* `--top <number>`, `-t <number>`
  Limits the output to the top X counter champions.

* `--all`, `-a`
  Includes champions with low pick rate or small sample size.

### Example

```bash
lolzy counter darius -t 5 -rk gold
```

This command shows the **top 5 counters against Darius** based on **gold-ranked** matches.

---

## Notes

* lolzy is designed to be fast and easy to use directly from the terminal.
* All data is based on the current League of Legends patch.
* Results may vary depending on rank and sample size.

---

Happy climbing! 🚀
