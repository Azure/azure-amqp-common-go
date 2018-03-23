# Azure AMQP Common
This project contains reusable components for AMQP based services like Event Hub and Service Bus. You will find 
abstractions over authentication, claims-based security, connection string parsing, checkpointing and RPC for AMQP.

If you are looking for the Azure Event Hub library for go, you can find it [here](https://github.com/Azure/azure-event-hubs-go).

## Install

### Via dep

```
dep ensure -add "github.com/Azure/azure-amqp-common-go"
```

### Or via go get
```
go get github.com/Azure/azure-amqp-common-go
```

## Contributing

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

## License

MIT, see [LICENSE](./LICENSE).

## Contribute

See [CONTRIBUTING.md](./contributing.md).