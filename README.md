Go Pool
================================

Not intended for production use - learning project. Sample usage:

```go
  import "pool"

  type Database struct {
    connections *pool.Pool
  }

  func (db *Database) GetConnection() (*Connection, error) {
    obj, e := db.connections.Get(20 * time.Second)

    if obj == nil {
      if e == ErrNoMember {
        return db.connections.NewConnection()
      } else {
        return nil, e
      }
    } else {
      return nil, e
    } 

    return m.(*Connection), nil
  }

  func (db *Database) PutConnection(connection *Connection) (error) {
    return db.connections.Put(connection)
  }

  func (db *Database) NewConnection() (*Connection, error) {
    connection := &Connection{}
    error = db.connections.Register(connection)
    return connection, error
  }

  func main(){
    db = &Database{connections: NewPool(20)}

    connection := db.GetConnection()

    ///////////////////////////////////////////////
    // Do something with the connection
    ///////////////////////////////////////////////

    db.PutConnection(connection)
  }
```
