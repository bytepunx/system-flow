# Blockquotes

TH-0101's entry quotes two steps of a list, each on its own.

> 3. Plan what your item's state calls for.

The rest of step 3 is unchanged. Step 6 now reads:

> 6. Your final message is the summary flai serve logs for your run.

The two steps in one quote, a quoted blank line between them.

> 3. Step three text
>
> 6. Step six text

A quoted list numbered in order.

> 1. One
> 2. Two
> 3. Three

A quoted list numbered one.

> 1. One
> 1. One again
> 1. One more

A quoted bullet list with a marker that changes.

> - Dash
> * Star

A list nested in a quoted list.

> 1. Outer
>    - Inner at three
>     - Inner at four
> 2. Outer again
>
>    5. Nested ordered from five

A quoted bullet list indented one.

>  - Indented one
>  - Indented one again

A nested quote with a list.

> > 2. Two
> > 5. Five

A nested quote that drops a level.

> > Inner paragraph
> 7. Lazy or not

A lazy continuation line.

> Quoted text
continued lazily with ** bad strong **.

Text directly before a quote.
> Quoted straight after the paragraph, with ** bad strong **.

A quote that ends at a blank line with a list after it.

> 1. Quoted

4. Unquoted from four

A quoted fence without a language.

> ```
> code	with a tab
> ```

A quoted fence with no blank line around it.

> Text
> ```text
> code
> ```
> Text again

A quoted heading.

> ## Quoted heading
> Text under it
> #### Skipped level

A quoted setext heading.

> Setext text
> ---

A quoted heading with trailing punctuation.

> ### Heading ends with a colon:

Quoted emphasis used as a heading.

> **Emphasis alone**

Quoted code with spaces.

> Run ` x ` and `y ` here.

A quoted bare URL.

> See https://example.com for more.

A quote without the optional space.

>1. Tight one
>3. Tight three

A quote in a list item.

- Item

  > 2. Quoted in an item
  > 4. Again
  >
  > - Quoted bullet in an item

A list item in a quote that runs on lazily.

> - Item text
continued lazily
> - Next item

Indented code in a quote.

>     indented code
>     more code

A quoted table.

> | a | b |
> | - | - |
> | c | *d * |

A quote opened on a list item's line.

- > Quoted on the item line
  > 3. Text, not an item

A nested quote with a nested list that drops a level.

> > - Deep
> >   - Deeper
>  - Back out one

A quoted comment before a list.

> <!-- comment
> still comment
> -->
> 2. After the comment

A quote with a tab after its marker.

>	- Tab after the marker
>	- Again

The end.
