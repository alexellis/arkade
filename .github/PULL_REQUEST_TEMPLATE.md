<!--- Provide a general summary of your changes in the Title above -->

## Description
<!--- Describe your changes in detail -->

## Motivation and Context
<!--- Why is this change required? What problem does it solve? -->
<!--- If it fixes an open issue, please link to the issue here. -->
- [ ] I have raised an issue to propose this change, which has been given a label of `design/approved` by a maintainer ([required](https://github.com/alexellis/arkade/blob/master/CONTRIBUTING.md))


## How Has This Been Tested?
<!--- Please describe in detail how you tested your changes. -->
<!--- Include details of your testing environment, and the tests you ran to -->
<!--- see how your change affects other areas of the code, etc. -->

If updating or adding a new CLI to `arkade get`, run:

```
go build && ./hack/test-tool.sh TOOL_NAME
```

## Types of changes
<!--- What types of changes does your code introduce? Put an `x` in all the boxes that apply: -->
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to change)

## Documentation

- [ ] I have updated the list of tools in README.md if (required) with `./arkade get --format markdown`
- [ ] I have updated the list of apps in README.md if (required) with `./arkade install --help`

## Checklist:
- [ ] My code follows the code style of this project.
- [ ] My change requires a change to the documentation.
- [ ] I have updated the documentation accordingly.
- [ ] I've read the [CONTRIBUTION](https://github.com/alexellis/arkade/blob/master/CONTRIBUTING.md) guide
- [ ] I have signed-off my commits with `git commit -s`

<!--- For "arkade install" or "arkade system install" -->
- [ ] I have tested this on arm, or have added code to prevent deployment
