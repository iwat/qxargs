qxargs
======

Call `xargs` quickly.

Well, the name may be over claimed because I did not implement all functionalities of `xargs`.

Usage
-----

This should be understandable.

```
SYNOPSIS
     qxargs [flags] command filters ...
     qxargs [flags] command commandargs ... -- filters ...

DESCRIPTION
     Execute command on the list of files that match the given filter.

     command  The command to be executed with the matches files.

     filters  There are 2 kinds of filters are supported, name filter and content filter.
              Simple string will be treated as file name filter.
              String with leading '?' will be treat as content filter.
              Multiple filters will be treated as AND.

EXAMPLES
     To execute vim on any file that has go in their name.

         $ qxargs  vim go

     To execute vim -p on any file that has go in their name.

         $ qxargs  vim -p -- go

     To execute vim -p on any file that has go in their name and has newGrepper in their contents.

         $ qxargs  vim -p -- go ?newGrepper
```

Animation?
----------

Sure, I love animation.

`qxargs vi fli`

![flow1](https://cloud.githubusercontent.com/assets/245383/25738837/6f011caa-31a9-11e7-8c00-e9d519a843a0.gif)

`qxargs vi -o -- find`

![flow2](https://cloud.githubusercontent.com/assets/245383/25738840/70f066ec-31a9-11e7-879f-24105ffec473.gif)

Interactive Mode
----------------

I have interactive mode too. Actually, there is no way to disable it right now.

![interactive](https://cloud.githubusercontent.com/assets/245383/25738931/d77e3470-31a9-11e7-9cc4-b873beb68e55.png)

License
-------

It's [MIT](https://github.com/iwat/qxargs/blob/master/LICENSE) license.

In short:

- Permissions
  - Commercial use
  - Modification
  - Distribution
  - Private use

- Conditions
  - License and copyright notice

- Limitations
  - Liability
  - Warranty
