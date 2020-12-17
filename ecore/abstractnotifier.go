package ecore

type notifierAdapterList struct {
	*basicEList
	notifier *AbstractNotifier
}

type notifierNotification struct {
	*abstractNotification
	notifier *AbstractNotifier
}

func (n *notifierNotification) GetNotifier() ENotifier {
	return n.notifier.interfaces.(ENotifier)
}

func (n *notifierNotification) GetFeature() EStructuralFeature {
	return nil
}

func (n *notifierNotification) GetFeatureID() int {
	return -1
}

func newNotifierNotification(notifier *AbstractNotifier, eventType EventType, oldValue interface{}, newValue interface{}, position int) *notifierNotification {
	n := new(notifierNotification)
	n.abstractNotification = NewAbstractNotification(eventType, oldValue, newValue, position)
	n.notifier = notifier
	return n
}

func newNotifierAdapterList(notifier *AbstractNotifier) *notifierAdapterList {
	l := new(notifierAdapterList)
	l.basicEList = NewEmptyBasicEList()
	l.notifier = notifier
	l.interfaces = l
	return l
}

func (l *notifierAdapterList) didAdd(index int, elem interface{}) {
	notifier := l.notifier.interfaces.(ENotifier)
	elem.(EAdapter).SetTarget(notifier)
}

func (l *notifierAdapterList) didRemove(index int, elem interface{}) {
	notifier := l.notifier.interfaces.(ENotifier)
	adapter := elem.(EAdapter)
	if notifier.EDeliver() {
		adapter.NotifyChanged(newNotifierNotification(l.notifier, REMOVING_ADAPTER, elem, nil, index))
	}
	adapter.UnSetTarget(notifier)
}

type AbstractNotifier struct {
	interfaces interface{}
}

type ENotifierInternal interface {
	ENotifier
	HasAdapters() bool
}

func NewAbstractNotifier() *AbstractNotifier {
	notifier := new(AbstractNotifier)
	notifier.interfaces = notifier
	return notifier
}

func (notifier *AbstractNotifier) AsENotifier() ENotifier {
	return notifier.interfaces.(ENotifier)
}

// SetInterfaces ...
func (notifier *AbstractNotifier) SetInterfaces(interfaces interface{}) {
	notifier.interfaces = interfaces
}

// GetInterfaces ...
func (notifier *AbstractNotifier) GetInterfaces() interface{} {
	return notifier.interfaces
}

func (notifier *AbstractNotifier) HasAdapters() bool {
	adapters := notifier.AsENotifier().EAdapters()
	return adapters != nil && !adapters.Empty()
}

func (notifier *AbstractNotifier) EAdapters() EList {
	return NewEmptyImmutableEList()
}

func (notifier *AbstractNotifier) EDeliver() bool {
	return false
}

func (notifier *AbstractNotifier) ESetDeliver(value bool) {
	panic("operation not supported")
}

func (notifier *AbstractNotifier) ENotify(notification ENotification) {
	n := notifier.AsENotifier()
	adapters := n.EAdapters()
	deliver := n.EDeliver()
	if adapters != nil && deliver {
		for it := adapters.Iterator(); it.HasNext(); {
			it.Next().(EAdapter).NotifyChanged(notification)
		}
	}
}

func (notifier *AbstractNotifier) ENotificationRequired() bool {
	n := notifier.interfaces.(ENotifierInternal)
	return n.HasAdapters() && n.EDeliver()
}
