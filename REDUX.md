# Redux

To implement redux we need to support Reduce operations that map from Action to State

func (o ObservableAction) Reduce(initial State, reducer func(State,Action)State) ObservableState

func (o ObservableAction) Scan(initial State, reducer func(State,Action)State) ObservableState

Reduce = Scan().Last()