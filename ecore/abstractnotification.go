package ecore

type abstractNotification struct {
	interfaces interface{}
	eventType  EventType
	oldValue   interface{}
	newValue   interface{}
	position   int
	next       ENotificationChain
}

func (notif *abstractNotification) GetEventType() EventType {
	return notif.eventType
}

func (notif *abstractNotification) GetOldValue() interface{} {
	return notif.oldValue
}

func (notif *abstractNotification) GetNewValue() interface{} {
	return notif.newValue
}

func (notif *abstractNotification) GetPosition() int {
	return notif.position
}

func (notif *abstractNotification) Merge(eOther ENotification) bool {
	eNotif := notif.interfaces.(ENotification)
	switch ev := notif.eventType; ev {
	case SET, UNSET:
		switch notifEv := eOther.GetEventType(); notifEv {
		case SET, UNSET:
			if eNotif.GetNotifier() == eOther.GetNotifier() &&
				eNotif.GetFeatureID() == eOther.GetFeatureID() {
				notif.newValue = eOther.GetNewValue()
				if eOther.GetEventType() == SET {
					notif.eventType = SET
				}
				return true
			}
		}
	case REMOVE:
		switch notifEv := eOther.GetEventType(); notifEv {
		case REMOVE:
			if eNotif.GetNotifier() == eOther.GetNotifier() &&
				eNotif.GetFeatureID() == eOther.GetFeatureID() {
				originalPosition := notif.GetPosition()
				notificationPosition := eOther.GetPosition()
				notif.eventType = REMOVE_MANY
				var removedValues []interface{}
				if originalPosition <= notificationPosition {
					removedValues = []interface{}{notif.oldValue, eOther.GetOldValue()}
					notif.position = originalPosition
					notif.newValue = []interface{}{originalPosition, notificationPosition + 1}
				} else {
					removedValues = []interface{}{eOther.GetOldValue(), notif.oldValue}
					notif.position = notificationPosition
					notif.newValue = []interface{}{notificationPosition, originalPosition}
				}
				notif.oldValue = removedValues
				return true
			}
		}
	case REMOVE_MANY:
		switch notifEv := eOther.GetEventType(); notifEv {
		case REMOVE:
			if eNotif.GetNotifier() == eOther.GetNotifier() &&
				eNotif.GetFeatureID() == eOther.GetFeatureID() {
				notificationPosition := eOther.GetPosition()
				positions := notif.newValue.([]interface{})
				newPositions := []interface{}{}

				index := 0
				for index < len(positions) {
					oldPosition := positions[index]
					if oldPosition.(int) <= notificationPosition {
						newPositions = append(newPositions, oldPosition)
						index++
						notificationPosition++
					} else {
						break
					}
				}

				oldValue := notif.oldValue.([]interface{})

				oldValue = append(oldValue, nil)
				copy(oldValue[index+1:], oldValue[index:])
				oldValue[index] = eOther.GetOldValue()

				newPositions = append(newPositions, notificationPosition)
				index++
				for index < len(positions) {
					newPositions = append(newPositions, positions[index-1])
					index++
				}
				notif.oldValue = oldValue
				notif.newValue = newPositions
				return true
			}
		}
	}
	return false
}

func (notif *abstractNotification) Add(eOther ENotification) bool {
	if eOther == nil {
		return false
	}
	if notif.Merge(eOther) {
		return false
	}
	if notif.next == nil {
		value, ok := eOther.(ENotificationChain)
		if ok {
			notif.next = value
			return true
		} else {
			notif.next = NewNotificationChain()
			return notif.next.Add(eOther)
		}
	} else {
		return notif.next.Add(eOther)
	}
}

func (notif *abstractNotification) Dispatch() {
	notification := notif.interfaces.(ENotification)
	notifier := notification.GetNotifier()
	if notifier != nil && notif.eventType != -1 {
		notifier.ENotify(notification)
	}
	if notif.next != nil {
		notif.next.Dispatch()
	}
}
