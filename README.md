# mcdeathregexp

This is a simple bit of go that will extract lang files from minecraft jar files and convert them to regexps.

Originally this was built for my bot GoGoGameBot to detect deaths on stdout on server versions I haven't modified.

## Usage

```
mcdeathregexp -mcjar ./some.jar
```

Additional options exist for simply dumping the lang entries as strings, ignoring keys, 
setting the path to the lang file in the jar