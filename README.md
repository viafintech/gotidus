gotidus
=======

[![Travis Build state](https://travis-ci.org/Barzahlen/gotidus.svg?branch=master)](https://travis-ci.org/Barzahlen/gotidus) [![GoDoc](https://godoc.org/github.com/Barzahlen/gotidus?status.png)](https://godoc.org/github.com/Barzahlen/gotidus)

gotidus is a Golang library which allows automatic view generation for every table in an SQL database. The purpose of the views is to anonymize the contents of select columns to ensure that no confidential information leaves the database while still providing access to the data in general.

It is also a port from the [tidus](https://github.com/Barzahlen/tidus) Ruby Gem.

## Install
```bash
go get github.com/Barzahlen/gotidus
```

## Usage

Please see the example below as well as the [godoc](https://godoc.org/github.com/Barzahlen/gotidus) reference.

### Example
```go
fooTable := gotidus.NewTable()
// Define columns on the table to anonymize in a specific way.
// Other columns will just contain their normal value.
// Note: Any column defined but not actually in the table will be ignored.
fooTable.AddAnonymizer(
    "bar",
    gotidus.NewStaticAnonymizer("staticValue", "TEXT"),
)
generator := gotidus.NewGenerator(postgres.NewQueryBuilder())
// Define tables that should have specifically anonymized columns.
// Tables that are not supposed to be anonymized specifically,
// do not have to be defined.
// 
// Note: Any table defined but not actually in the database will be ignored.
generator.AddTable("foo", fooTable)

// Clear existing views
err := generator.ClearViews(db)
if err != nil {
    log.Fatal(err)
}

// ... database migration

// Create new views
err = generator.CreateViews(db)
if err != nil {
    log.Fatal(err)
}
```

## Backup and Restore

You can use the bash example script located in examples to backup and restore databases prepared with tidus easily. `tidus_backup_restore.sh` can be called with any parameter other than `-d|-r|--dump|--restore` to get help for it's usage. The `tidus_seq_rst.sql` file is necessary for restores since it's will reset all sequences after restore for you - it's not necessary for backups only.
You also need the `tidus_credentials.conf` with the IP/DNS, User and Password of the Dump and Restore users. If you use `tidus_backup_restore.sh` on separate machines for backup and restore, you can split up the credentials file and only provide the information necessary to backup and restore.

### Basic usage

Before dumping or restoring you have to provide the `tidus_credentials.conf` file with all the informations needed for dumping and restoring. Those parameters are not exposed into the commandline due to security considerations. Also manually edit the `tidus_backup_restore.sh` and check the `dump_it` and `restore_it` functions and add the databases you want to dump or restore as well as the database names in your staging environment and the staging user which will get the permissions after restore.

- `./tidus_backup_restore.sh /path/to/tidus_credentials.conf -d /path/to/the/dumps/folder`
  - Add all databases you want to dump from in the `dump_it` function!
- `./tidus_backup_restore.sh /path/to/tidus_credentials.conf -r /path/to/the/dumps/folder <Backup-Set-No>`
  - Add all databases you want to restore - as well as the destination database names and users - in the `restore_it` function!
  - Be sure to have the `tidus_seq_rst.sql`in the same folder as the script which is required for a successful restore!

## Bugs and Contribution
For bugs and feature requests open an issue on Github. For code contributions fork the repo, make your changes and create a pull request.

## Extending functionality
The number of anonymizers implemented so far is limited.
A new anonymization strategy can be easily defined through implementation of the `gotidus.Anonymizer` interface.
It is furthermore possible to add support for other databases by implementing the `gotidus.QueryBuilder` interface.

## License
[LICENSE](LICENSE)
