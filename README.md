[![Build Status](https://dev.azure.com/BladeMonRTPipelines/BladeMonRT/_apis/build/status/microsoft.BladeMonRT?branchName=main)](https://dev.azure.com/BladeMonRTPipelines/BladeMonRT/_build/latest?definitionId=1&branchName=main)

# Project

> This repo has been populated by an initial template to help get you started. Please
> make sure to update the content to build a great experience for community-building.

As the maintainer of this project, please make a few updates:

- Improving this README.MD file to provide a great experience
- Updating SUPPORT.MD with content about this project's support experience
- Understanding the security reporting process in SECURITY.MD
- Remove this section from the README

## Contributing

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.opensource.microsoft.com.

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

#### Submitting a Pull Request
- This project follows the GitHubFlow branching technique and all branches should be named with the *feature/* prefix.
- A new Azure DevOps build is created for every commit to a feature branch or to the master. 
- GitVersion is used to automate package versioning.
- feature branches have the *alpha* tag appended.

## Trademarks

This project may contain trademarks or logos for projects, products, or services. Authorized use of Microsoft 
trademarks or logos is subject to and must follow 
[Microsoft's Trademark & Brand Guidelines](https://www.microsoft.com/en-us/legal/intellectualproperty/trademarks/usage/general).
Use of Microsoft trademarks or logos in modified versions of this project must not cause confusion or imply Microsoft sponsorship.
Any use of third-party trademarks or logos are subject to those third-party's policies.

## Run and Test
#### To run the main use:

cd BladeMonRT

./run.bat

#### To run all tests except the end-to-end test use:

cd BladeMonRT

./test.bat

#### To run the end-to-end test:

go test -run TestEndToEnd

* The end-to-end test runs BRT until a keyboard interrupt.
* You will have to manually raise ETW events in a separate terminal.

To create a mock of a GO interface

1. Use the mockgen command described here https://github.com/golang/mock

    Example: 

    cd BladeMonRT

    mockgen -source="./nodes/node.go" -destination="./nodes/mock_node.go" -package="nodes"

2. Then in the GO file that contains the interface mocked, add a comment following this
template after the imports:

// [MockedInterface] mock generation.

//go:generate [mockgen command]

Example:

// InterfaceNodeFactory mock generation.

//go:generate mockgen -source=./node_factory.go -destination=./mock_node_factory.go -package=main

This will ensure when test.bat is run, the mock for the [MockedInterface] is regenerated.

#### To format all GO files run:

cd BladeMonRT

gofmt -l -s -w .
