# test

    import "github.com/reactivego/rx/test"

Package `test` provides tests for the generic rx package.

Every operator / data type has its own subdirectory named after it containing one or more tests that exercise its functionality.

## Operators
Below is the list of implemented [operators](http://reactivex.io/documentation/operators.html). Operators with a :star: are the most commonly used ones.

### Creating Operators
Operators that originate new Observables.

- [**Create**](https://godoc.org/github.com/reactivego/rx/test/Create/)() :star: Observable
- [**CreateRecursive**](https://godoc.org/github.com/reactivego/rx/test/CreateRecursive/)() Observable
- [**CreateFutureRecursive**](https://godoc.org/github.com/reactivego/rx/test/CreateFutureRecursive/)() Observable
- [**Defer**](https://godoc.org/github.com/reactivego/rx/test/Defer/)() Observable
- [**Empty**](https://godoc.org/github.com/reactivego/rx/test/Empty/)() Observable
- [**From**](https://godoc.org/github.com/reactivego/rx/test/From/)() :star: Observable
- [**FromChan**](https://godoc.org/github.com/reactivego/rx/test/FromChan/)() Observable
- [**Just**](https://godoc.org/github.com/reactivego/rx/test/Just/)() :star: Observable
- [**Never**](https://godoc.org/github.com/reactivego/rx/test/Never/)() Observable
- [**Of**](https://godoc.org/github.com/reactivego/rx/test/Of/)() :star: Observable
- [**Range**](https://godoc.org/github.com/reactivego/rx/test/Range/)() ObservableInt
- [**Start**](https://godoc.org/github.com/reactivego/rx/test/Start/)() Observable
- [**Throw**](https://godoc.org/github.com/reactivego/rx/test/Throw/)() Observable

### Transforming Operators
Operators that transform items that are emitted by an Observable.

- (Observable) [**Map**](https://godoc.org/github.com/reactivego/rx/test/Map/)() :star: Observable
- (Observable) [**Scan**](https://godoc.org/github.com/reactivego/rx/test/Scan/)() :star: Observable

### Filtering Operators
Operators that selectively emit items from a source Observable.

- (Observable) [**Distinct**](https://godoc.org/github.com/reactivego/rx/test/Distinct/)() Observable
- (Observable) [**ElementAt**](https://godoc.org/github.com/reactivego/rx/test/ElementAt/)() Observable
- (Observable) [**Filter**](https://godoc.org/github.com/reactivego/rx/test/Filter/)() :star: Observable
- (Observable) [**First**](https://godoc.org/github.com/reactivego/rx/test/First/)() Observable
- (Observable) [**IgnoreElements**](https://godoc.org/github.com/reactivego/rx/test/IgnoreElements/)() Observable
- (Observable) [**IgnoreCompletion**](https://godoc.org/github.com/reactivego/rx/test/IgnoreCompletion/)() Observable
- (Observable) [**Last**](https://godoc.org/github.com/reactivego/rx/test/Last/)() Observable
- (Observable) [**Single**](https://godoc.org/github.com/reactivego/rx/test/Single/)() Observable
- (Observable) [**Skip**](https://godoc.org/github.com/reactivego/rx/test/Skip/)() Observable
- (Observable) [**SkipLast**](https://godoc.org/github.com/reactivego/rx/test/SkipLast/)() Observable
- (Observable) [**Take**](https://godoc.org/github.com/reactivego/rx/test/Take/)() :star: Observable
- (Observable) [**TakeLast**](https://godoc.org/github.com/reactivego/rx/test/TakeLast/)() Observable
- (Observable) [**TakeUntil**](https://godoc.org/github.com/reactivego/rx/test/TakeUntil/)() :star: Observable
- (Observable) [**TakeWhile**](https://godoc.org/github.com/reactivego/rx/test/TakeWhile/)() :star: Observable

### Time Based Operators

Time based creating and filtering operators.

- (Observable) [**Audit**](https://godoc.org/github.com/reactivego/rx/test/Audit/)() Observable
- (Observable) [**Debounce**](https://godoc.org/github.com/reactivego/rx/test/Debounce/)() :star: Observable
- (Observable) [**Delay**](https://godoc.org/github.com/reactivego/rx/test/Delay/)() Observable
- [**Interval**](https://godoc.org/github.com/reactivego/rx/test/Interval/)() ObservableInt
- (Observable) [**Sample**](https://godoc.org/github.com/reactivego/rx/test/Sample/)() Observable
- (Observable) [**Throttle**](https://godoc.org/github.com/reactivego/rx/test/Throttle/)() Observable
- [**Ticker**](https://godoc.org/github.com/reactivego/rx/test/Ticker/)() ObservableTime
- (Observable) [**Timeout**](https://godoc.org/github.com/reactivego/rx/test/Timeout/)() Observable
- [**Timer**](https://godoc.org/github.com/reactivego/rx/test/Timer/)() ObservableInt
- (Observable) [**Timestamp**](https://godoc.org/github.com/reactivego/rx/test/Timestamp/)() Observable
- (Observable) [**TimeInterval**](https://godoc.org/github.com/reactivego/rx/test/TimeInterval/)() Observable

### Combining Operators
Operators that work with multiple source observables to create a single observable. It looks like there is an underlying logic at play for naming the different kinds of combining operators. The matrices below guides the naming of the operators. Where operators don't make sense we've placed a `-` in the cell. If it is not known yet whether an operator makes sense, a `?` is placed in the cell.

**Concat**, **Switch** and **Exhaust** are all operators that flatten the emissions of multiple observables into a single stream by subscribing to every observable *stricly in sequence*. Observables may be added while the flattening is already going on.

| Function            | Operator                | Map                    | MapTo             | All             |
|:--------------------|:------------------------|:-----------------------|:------------------|:----------------|
| [**Concat**] :star: | [**ConcatWith**] :star: | [**ConcatMap**]        | [**ConcatMapTo**] | [**ConcatAll**] |
| -                   | -                       | [**SwitchMap**] :star: | SwitchMapTo       | [**SwitchAll**] |
| -                   | -                       | ExhaustMap             | ExhaustMapTo      | ExhaustAll      |

[**Concat**]: https://godoc.org/github.com/reactivego/rx/test/Concat/
[**ConcatWith**]: https://godoc.org/github.com/reactivego/rx/test/ConcatWith/
[**ConcatMap**]: https://godoc.org/github.com/reactivego/rx/test/ConcatMap/
[**ConcatMapTo**]: https://godoc.org/github.com/reactivego/rx/test/ConcatMapTo/
[**ConcatAll**]: https://godoc.org/github.com/reactivego/rx/test/ConcatAll/
[**SwitchMap**]: https://godoc.org/github.com/reactivego/rx/test/SwitchMap/
[**SwitchAll**]: https://godoc.org/github.com/reactivego/rx/test/SwitchAll/

**Merge**, **MergeDelayError** and **Race** are operators that flatten the emissions of multiple observables into a single stream by subscribing to all observables *concurrently*. Here also, observables may be added while the flattening is already going on.

| Function              | Operator                  | Map                   | MapTo                | All                |
|:----------------------|:--------------------------|:----------------------|:---------------------|:-------------------|
| [**Merge**] :star:    | [**MergeWith**] :star:    | [**MergeMap**] :star: | MergeMapTo           | [**MergeAll**]     |
| [**MergeDelayError**] | [**MergeDelayErrorWith**] | MergeDelayErrorMap    | MergeDelayErrorMapTo | MergeDelayErrorAll |
| Race                  | RaceWith                  | RaceMap               | RaceMapTo            | RaceAll            |

[**Merge**]: https://godoc.org/github.com/reactivego/rx/test/Merge/
[**MergeWith**]: https://godoc.org/github.com/reactivego/rx/test/MergeWith/
[**MergeMap**]: https://godoc.org/github.com/reactivego/rx/test/MergeMap/
[**MergeAll**]: https://godoc.org/github.com/reactivego/rx/test/MergeAll/
[**MergeDelayError**]: https://godoc.org/github.com/reactivego/rx/test/MergeDelayError/
[**MergeDelayErrorWith**]: https://godoc.org/github.com/reactivego/rx/test/MergeDelayErrorWith/


**CombineLatest**, **Zip**, **WithLatestFrom** and **ForkJoin** are operators that flatten the emissions of multiple observables into a single observable that emits slices of values. Differently from the previous two sets of operators, these operators only start emitting once the list of observables to flatten is complete. 

| Function            | Operator                | Map                    | MapTo                    | All                    |
|:--------------------|:------------------------|:-----------------------|:-------------------------|:-----------------------|
| [**CombineLatest**] | [**CombineLatestWith**] | [**CombineLatestMap**] | [**CombineLatestMapTo**] | [**CombineLatestAll**] |
| Zip                 | ZipWith                 | ZipMap                 | ZipMapTo                 | ZipAll                 |
| ?                   | WithLatestFrom :star:   | ?                      | ?                        | ?                      |
| ForkJoin            | ?                       | ?                      | ?                        | ?                      |

[**CombineLatest**]: https://godoc.org/github.com/reactivego/rx/test/CombineLatest/
[**CombineLatestWith**]: https://godoc.org/github.com/reactivego/rx/test/CombineLatestWith/
[**CombineLatestMap**]: https://godoc.org/github.com/reactivego/rx/test/CombineLatestMap/
[**CombineLatestMapTo**]: https://godoc.org/github.com/reactivego/rx/test/CombineLatestMapTo/
[**CombineLatestAll**]: https://godoc.org/github.com/reactivego/rx/test/CombineLatestAll/

### Error Handling Operators
Operators that help to recover from error notifications from an Observable.

- (Observable) [**Catch**](https://godoc.org/github.com/reactivego/rx/test/Catch/)() :star: Observable
- (Observable) [**Retry**](https://godoc.org/github.com/reactivego/rx/test/Retry/)() Observable

### Utility Operators
A toolbox of useful Operators for working with Observables.

- (Observable) [**Do**](https://godoc.org/github.com/reactivego/rx/test/Do/)() :star: Observable
- (Observable) [**DoOnError**](https://godoc.org/github.com/reactivego/rx/test/DoOnError/)() Observable
- (Observable) [**DoOnComplete**](https://godoc.org/github.com/reactivego/rx/test/DoOnComplete/)() Observable
- (Observable) [**Finally**](https://godoc.org/github.com/reactivego/rx/test/Finally/)() Observable
- (Observable) [**Passthrough**](https://godoc.org/github.com/reactivego/rx/test/Passthrough/)() Observable
- (Observable) [**Repeat**](https://godoc.org/github.com/reactivego/rx/test/Repeat/)() Observable
- (Observable) [**Serialize**](https://godoc.org/github.com/reactivego/rx/test/Serialize/)() Observable

### Conditional and Boolean Operators
Operators that evaluate one or more Observables or items emitted by Observables.

- (Observable) [**All**](https://godoc.org/github.com/reactivego/rx/test/All/)() ObservableBool

### Aggregate Operators
Operators that operate on the entire sequence of items emitted by an Observable.

- (Observable) [**Average**](https://godoc.org/github.com/reactivego/rx/test/Average/)() Observable
- (Observable) [**Count**](https://godoc.org/github.com/reactivego/rx/test/Count/)() ObservableInt
- (Observable) [**Max**](https://godoc.org/github.com/reactivego/rx/test/Max/)() Observable
- (Observable) [**Min**](https://godoc.org/github.com/reactivego/rx/test/Min/)() Observable
- (Observable) [**Reduce**](https://godoc.org/github.com/reactivego/rx/test/Reduce/)() Observable
- (Observable) [**Sum**](https://godoc.org/github.com/reactivego/rx/test/Sum/)() Observable

### Type Casting and Type Filtering Operators
Operators to type cast, type filter observables.

- (Observable) [**AsObservable**](https://godoc.org/github.com/reactivego/rx/test/AsObservable)() Observable
- (Observable) [**Only**](https://godoc.org/github.com/reactivego/rx/test/Only/)() Observable

### Scheduling Operators
Change the scheduler for subscribing and observing.

- (Observable) [**ObserveOn**](https://godoc.org/github.com/reactivego/rx/test/ObserveOn/)() Observable
- (Observable) [**SubscribeOn**](https://godoc.org/github.com/reactivego/rx/test/SubscribeOn/)() Observable

### Multicasting Operators
A *Multicaster* is a combination of an *Observable* and a *Connectable* that will multicast to the observers that subscribe to it.
The *Multicaster* itself will subscribe to its *Observable* when the *Connect* method is called on it.

- (Observable) [**Publish**](https://godoc.org/github.com/reactivego/rx/test/Publish/)() Multicaster
- (Observable) [**PublishReplay**](https://godoc.org/github.com/reactivego/rx/test/PublishReplay/)() Multicaster

*Multicaster* supports different strategies for subscribing to the *Observable* from which it was created.

- (Multicaster) [**RefCount**](https://godoc.org/github.com/reactivego/rx/test/RefCount/)() Observable
- (Multicaster) [**AutoConnect**](https://godoc.org/github.com/reactivego/rx/test/AutoConnect/)() Observable

*Connectable* is a data type that implements the Connect method. It is used as part of a *Multicaster*.
- (Connectable) [**Connect**](https://godoc.org/github.com/reactivego/rx/test/Connect)() Subscription

## Subjects
A *Subject* is both a multicasting *Observable* as well as an *Observer*. The *Observable* side allows multiple simultaneous subscribers. The *Observer* side allows you to directly feed it data or subscribe it to another *Observable*.

- [**Subject**](https://godoc.org/github.com/reactivego/rx/test/Subject)() Subject
- [**ReplaySubject**](https://godoc.org/github.com/reactivego/rx/test/ReplaySubject)() Subject

## Subscribing
Subscribing breathes life into a chain of observables. An observable may be subscribed to many times. 

`type Observable`
- (Observable) [**Println**](https://godoc.org/github.com/reactivego/rx/test/Println)() error
- (Observable) [**Subscribe**](https://godoc.org/github.com/reactivego/rx/test/Subscribe)() Subscription
- (Observable) [**ToChan**](https://godoc.org/github.com/reactivego/rx/test/ToChan)() chan foo
- (Observable) [**ToSingle**](https://godoc.org/github.com/reactivego/rx/test/ToSingle)() (foo, error)
- (Observable) [**ToSlice**](https://godoc.org/github.com/reactivego/rx/test/ToSlice)() ([]foo, error)
- (Observable) [**Wait**](https://godoc.org/github.com/reactivego/rx/test/Wait)() error

`type Multicaster struct{ Observable; Connectable }`
- (Connectable) [**Connect**](https://godoc.org/github.com/reactivego/rx/test/Connect)() Subscription
- (Multicaster) [**AutoConnect**](https://godoc.org/github.com/reactivego/rx/test/AutoConnect)() Observable
- (Multicaster) [**RefCount**](https://godoc.org/github.com/reactivego/rx/test/RefCount)() Observable
> *RefCount* and *AutoConnect* call *Connect*.

