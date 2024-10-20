# M'key

*Make working with multi-field keys in AWS Dynamodb simple*

AWS Dynamodb is a powerful, fully managed No SQL database solution, that can scale incredibly and cost-effective.
However, it requires developers to use different techniques and manage 'app side' some concerns that would traditionally
be handled by the database itself.

## Data Integrity in a Schema-less Space

While **Dynamodb** is itself schema-less, your data probably isn't.

See this best practice guide from
AWS: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/bp-sort-keys.html for the basics.