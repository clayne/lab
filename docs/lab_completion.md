## lab completion

Generates the shell autocompletion [bash, elvish, fish, oil, powershell, xonsh, zsh]

### Synopsis

Generates the shell autocompletion [bash, elvish, fish, oil, powershell, xonsh, zsh]

Scripts can be directly sourced (though using pre-generated versions is recommended to avoid shell startup delay):
  bash       : source <(lab completion)
  elvish     : eval(lab completion|slurp)
  fish       : lab completion | source
  oil        : source <(lab completion)
  powershell : lab completion | Out-String | Invoke-Expression
  xonsh      : exec($(lab completion))
  zsh        : source <(lab completion)

```
lab completion [shell] [flags]
```

### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
      --debug      Enable debug logging level
      --no-pager   Do not pipe output into a pager
      --quiet      Turn off any sort of logging. Only command output is printed
```

### SEE ALSO

* [lab](index.md)	 - lab: A GitLab Command Line Interface Utility

###### Auto generated by spf13/cobra on 20-May-2021
