package factory
import "owdile-notification/lib/impl"
import "owdile-notification/lib/base"


func  GetListener(li *base.ListenerInfo) base.IListener{
    switch li.Type {
    case "NormalListener":
        var n impl.NormalListener
        n.ListenerInfo = *li
        return n
    default:
        var n impl.NormalListener
        n.ListenerInfo = *li
        return n
    }

}


func  GetNotify(li *base.NotifyInfo) base.INotify{
    switch li.Type {
    case "NormalNotify":
           var n impl.NormalNotify
           n.NotifyInfo = *li
           return &n
    default:
        var n impl.NormalNotify
        n.NotifyInfo = *li
        return &n
    }

}
