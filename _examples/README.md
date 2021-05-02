## In depth examples

In this folder you'll find a `example.js` file containing few examples and a corresponding `.bobconfig.yml`

### How to use it

* You don't have any config file yet

```bash
$> git clone git@github.com:jcalmat/bob.git
$> cd bob
$> sed -i 's?/path\/to\/this\/project?'`pwd`'?' examples/.bobconfig.yml
$> cp examples/.bobconfig.yml ~/ #copy the provided config file
$> make install
$> bob build example
```

* You already have a config file

```bash
$> git clone git@github.com:jcalmat/bob.git
$> cd bob
$> sed -i 's?/path\/to\/this\/project?'`pwd`'?' examples/.bobconfig.yml
```

Take the command + template from the [config file](.bobconfig.yml) and add them to your own config file

```bash
$> make install
$> bob build example
```
