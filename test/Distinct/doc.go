/*
Distinct suppress duplicate items emitted by an Observable.

The operator filters an Observable by only allowing items through that have not already been emitted.
In some implementations there are variants that allow you to adjust the criteria by which two items are
considered “distinct.” In some, there is a variant of the operator that only compares an item against its
immediate predecessor for distinctness, thereby filtering only consecutive duplicate items from the sequence.
*/
package Distinct
