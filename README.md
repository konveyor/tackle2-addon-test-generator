[![Docker Repository on Quay](https://quay.io/repository/konveyor/tackle2-addon-test-generator/status "Docker Repository on Quay")](https://quay.io/repository/konveyor/tackle2-addon-test-generator)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkonveyor%2Ftackle2-addon-test-generator.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkonveyor%2Ftackle2-addon-test-generator?ref=badge_shield)

# Development of Tackle2 addon with test generator

The tool which should be implemented as an addon: https://github.com/konveyor/tackle-test-generator-cli

Tackle 2 addon example code: https://github.com/konveyor/tackle2-addon-windup and https://github.com/konveyor/tackle2-addon

## Steps of the addon integration

- addon fetch the application codebase locally
- addon creates toml config file for the application
- addon executes tkltest locally
- tkltest creates test files locally to the application codebase
- addon add, commit & push generated files into application origin location

## Development notes

Sample image: ```quay.io/maufart/tackle2-addon-test-generator```


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkonveyor%2Ftackle2-addon-test-generator.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkonveyor%2Ftackle2-addon-test-generator?ref=badge_large)