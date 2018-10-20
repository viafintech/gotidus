/*
  SQL anonymization view builder for go

  Example:

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
*/
package gotidus
