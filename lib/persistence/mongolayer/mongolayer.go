package mongolayer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
	moviepb "github.com/AkashGit21/ms-project/internal/grpc/movie"
	"github.com/AkashGit21/ms-project/lib/persistence"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

const (
	DATABASE = "msZone"
	USERNAME = "db_username"
	PASSWORD = "db_password"
	CLUSTER  = "db_cluster"
	USERS    = "users"
	MOVIES   = "movies"
)

type MongoDBLayer struct {
	client *mongo.Client
}

func NewMongoDBLayer(connection string) (persistence.DatabaseHandler, error) {
	var dbUsername, dbPassword, dbCluster string
	if dbUsername = os.Getenv("DB_USERNAME"); dbUsername == "" {
		dbUsername = USERNAME
	}
	if dbPassword = os.Getenv("DB_PASSWORD"); dbPassword == "" {
		dbPassword = PASSWORD
	}
	if dbCluster = os.Getenv("DB_CLUSTER"); dbCluster == "" {
		dbCluster = CLUSTER
	}

	uri, err := getConnectionURI(DATABASE, dbUsername, dbPassword, dbCluster)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("MONGODB URL:", uri)
	c, err := getNewClient(uri)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = c.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected!")

	return &MongoDBLayer{
		client: c,
	}, err
}

func (mgoLayer *MongoDBLayer) AddUser(u persistence.User) ([]byte, error) {
	cli := mgoLayer.client

	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	sess, err := cli.StartSession(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer sess.EndSession(context.TODO())
	var id []byte
	// Call WithSession to start a transaction within the new session.
	err = mongo.WithSession(
		context.TODO(),
		sess,
		func(sessCtx mongo.SessionContext) error {
			// Use sessCtx as the Context parameter for InsertOne and FindOne so
			// both operations are run under the new Session.

			if err := sess.StartTransaction(); err != nil {
				return err
			}

			usersCollection := cli.Database(DATABASE).Collection(USERS)
			res, err := usersCollection.InsertOne(sessCtx, u)
			if err != nil {
				log.Println(err)
				return err
			}

			id, _ = json.Marshal(res.InsertedID)

			return sess.CommitTransaction(context.Background())
		})

	return id, err
}

func (mgoLayer *MongoDBLayer) FindByUsername(uname string) (persistence.User, error) {
	cli := mgoLayer.client

	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())

	sess, err := cli.StartSession(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer sess.EndSession(context.TODO())

	var result persistence.User
	// Call WithSession to start a transaction within the new session.
	err = mongo.WithSession(
		context.TODO(),
		sess,
		func(sessCtx mongo.SessionContext) error {
			// Use sessCtx as the Context parameter for InsertOne and FindOne so
			// both operations are run under the new Session.

			if err := sess.StartTransaction(); err != nil {
				return err
			}

			newId := bson.M{"username": uname}
			// newId := bson.ObjectIdHex(id)
			usersCollection := cli.Database(DATABASE).Collection(USERS)
			err = usersCollection.FindOne(
				sessCtx,
				newId,
			).Decode(&result)
			if err != nil {
				log.Println(err)
				return err
			}

			return sess.CommitTransaction(context.Background())
		})

	return result, err
}

func (mgoLayer *MongoDBLayer) FindAllUsers(offset int, pgSize int32) ([]*identitypb.User, error) {
	cli := mgoLayer.client

	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())

	sess, err := cli.StartSession(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer sess.EndSession(context.TODO())

	var results []*identitypb.User
	// Call WithSession to start a transaction within the new session.
	err = mongo.WithSession(
		context.TODO(),
		sess,
		func(sessCtx mongo.SessionContext) error {
			// Use sessCtx as the Context parameter for InsertOne and FindOne so
			// both operations are run under the new Session.

			if err := sess.StartTransaction(); err != nil {
				return err
			}

			filter := bson.M{}
			opts := options.Find().SetSort(bson.M{}).SetSkip(int64(offset)).SetLimit(int64(pgSize))

			usersCollection := cli.Database(DATABASE).Collection(USERS)
			cur, err := usersCollection.Find(sessCtx, filter, opts)
			// defer cur.Close(ctx)

			if err = cur.All(sessCtx, &results); err != nil {
				log.Println(err)
				return err
			}
			return sess.CommitTransaction(context.Background())
		})

	// fmt.Println("ID0:", results[0].Username)
	return results, err
}

func (mgoLayer *MongoDBLayer) CountUsers() int {
	cli := mgoLayer.client

	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())

	sess, err := cli.StartSession(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer sess.EndSession(context.TODO())

	var count int64
	// Call WithSession to start a transaction within the new session.
	err = mongo.WithSession(
		context.TODO(),
		sess,
		func(sessCtx mongo.SessionContext) error {
			// Use sessCtx as the Context parameter for InsertOne and FindOne so
			// both operations are run under the new Session.

			if err := sess.StartTransaction(); err != nil {
				return err
			}

			filter := bson.M{}

			usersCollection := cli.Database(DATABASE).Collection(USERS)
			count, err = usersCollection.CountDocuments(sessCtx, filter)
			// defer cur.Close(ctx)

			if err != nil {
				log.Println(err)
				return err
			}
			return sess.CommitTransaction(context.Background())
		})

	// fmt.Println("ID0:", results[0].Username)
	return int(count)
}

func (mgoLayer *MongoDBLayer) RemoveByUsername(uname string) error {
	cli := mgoLayer.client

	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())

	sess, err := cli.StartSession(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer sess.EndSession(context.TODO())
	// var count []byte

	// Call WithSession to start a transaction within the new session.
	err = mongo.WithSession(
		context.TODO(),
		sess,
		func(sessCtx mongo.SessionContext) error {
			// Use sessCtx as the Context parameter for InsertOne and FindOne so
			// both operations are run under the new Session.

			if err := sess.StartTransaction(); err != nil {
				return err
			}

			newId := bson.M{"username": uname}
			// newId := bson.ObjectIdHex(id)
			usersCollection := cli.Database(DATABASE).Collection(USERS)
			_, err := usersCollection.DeleteOne(sessCtx, newId)
			if err != nil {
				log.Println(err)
				return err
			}

			// count, _ = json.Marshal(res.DeletedCount)

			return sess.CommitTransaction(context.Background())
		})

	return nil
}

func (mgoLayer *MongoDBLayer) Authenticate(uname string, password string) bool {
	cli := mgoLayer.client

	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())

	sess, err := cli.StartSession(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer sess.EndSession(context.TODO())

	var result persistence.User
	// Call WithSession to start a transaction within the new session.
	err = mongo.WithSession(
		context.TODO(),
		sess,
		func(sessCtx mongo.SessionContext) error {
			// Use sessCtx as the Context parameter for InsertOne and FindOne so
			// both operations are run under the new Session.

			if err := sess.StartTransaction(); err != nil {
				return err
			}

			filter := bson.M{"username": uname}
			// projection := bson.D{{"username", 1}, {"password", 0}}
			// options.Find().SetSort(bson.M{"username": uname}).SetProjection(bson.D{{"username", 1}, {"password", 0}})

			opts := options.FindOne().SetProjection(bson.M{"username": true, "password": true})

			usersCollection := cli.Database(DATABASE).Collection(USERS)
			usersCollection.FindOne(sessCtx, filter, opts).Decode(&result)
			if err != nil {
				return err
			}

			return sess.CommitTransaction(context.Background())
		})
	log.Println("Result: ", result)
	if result.Username == "" {
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password))
	if err != nil {
		log.Println("Password doesn't match! ", err)
		return false
	}

	return true
}

func (mgoLayer *MongoDBLayer) AddMovie(mv persistence.Movie) ([]byte, error) {
	cli := mgoLayer.client

	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	sess, err := cli.StartSession(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer sess.EndSession(context.TODO())
	var id []byte
	// Call WithSession to start a transaction within the new session.
	err = mongo.WithSession(
		context.TODO(),
		sess,
		func(sessCtx mongo.SessionContext) error {
			// Use sessCtx as the Context parameter for InsertOne and FindOne so
			// both operations are run under the new Session.

			if err := sess.StartTransaction(); err != nil {
				return err
			}

			moviesCollection := cli.Database(DATABASE).Collection(MOVIES)
			res, err := moviesCollection.InsertOne(sessCtx, mv)
			if err != nil {
				log.Println(err)
				return err
			}

			id, _ = json.Marshal(res.InsertedID)

			return sess.CommitTransaction(context.Background())
		})

	return id, err
}

func (mgoLayer *MongoDBLayer) FindMovieByID(id string) (persistence.Movie, error) {
	cli := mgoLayer.client

	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())

	sess, err := cli.StartSession(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer sess.EndSession(context.TODO())

	var result persistence.Movie
	// Call WithSession to start a transaction within the new session.
	err = mongo.WithSession(
		context.TODO(),
		sess,
		func(sessCtx mongo.SessionContext) error {
			// Use sessCtx as the Context parameter for InsertOne and FindOne so
			// both operations are run under the new Session.

			if err := sess.StartTransaction(); err != nil {
				return err
			}

			filter := bson.M{"_id": id}
			// newId := bson.ObjectIdHex(id)
			moviesCollection := cli.Database(DATABASE).Collection(MOVIES)
			err = moviesCollection.FindOne(
				sessCtx,
				filter,
			).Decode(&result)
			if err != nil {
				log.Println(err)
				return err
			}

			return sess.CommitTransaction(context.Background())
		})

	return result, err
}

func (mgoLayer *MongoDBLayer) FindAllMovies(offset int, pgSize int32) ([]*moviepb.Movie, error) {
	cli := mgoLayer.client

	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())

	sess, err := cli.StartSession(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer sess.EndSession(context.TODO())

	var results []*moviepb.Movie
	// Call WithSession to start a transaction within the new session.
	err = mongo.WithSession(
		context.TODO(),
		sess,
		func(sessCtx mongo.SessionContext) error {
			// Use sessCtx as the Context parameter for InsertOne and FindOne so
			// both operations are run under the new Session.

			if err := sess.StartTransaction(); err != nil {
				return err
			}

			filter := bson.M{}
			opts := options.Find().SetSort(bson.M{}).SetSkip(int64(offset)).SetLimit(int64(pgSize))

			moviesCollection := cli.Database(DATABASE).Collection(MOVIES)
			cur, err := moviesCollection.Find(sessCtx, filter, opts)
			// defer cur.Close(ctx)

			if err = cur.All(sessCtx, &results); err != nil {
				log.Println(err)
				return err
			}
			return sess.CommitTransaction(context.Background())
		})

	// fmt.Println("ID0:", results[0].Username)
	return results, err
}

func (mgoLayer *MongoDBLayer) UpdateMovieByID(id string, mv persistence.Movie) ([]byte, error) {
	cli := mgoLayer.client

	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	sess, err := cli.StartSession(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer sess.EndSession(context.TODO())

	var result persistence.Movie
	// Call WithSession to start a transaction within the new session.
	err = mongo.WithSession(
		context.TODO(),
		sess,
		func(sessCtx mongo.SessionContext) error {
			// Use sessCtx as the Context parameter for InsertOne and FindOne so
			// both operations are run under the new Session.

			if err := sess.StartTransaction(); err != nil {
				return err
			}

			filter := bson.M{"_id": id}
			moviesCollection := cli.Database(DATABASE).Collection(MOVIES)
			err := moviesCollection.FindOneAndUpdate(
				sessCtx,
				filter,
				mv).Decode(&result)
			if err != nil {
				log.Println(err)
				return err
			}
			// InsertOne(sessCtx, mv)
			if err != nil {
				log.Println(err)
				return err
			}

			return sess.CommitTransaction(context.Background())
		})

	return []byte(result.Id), err
}

func (mgoLayer *MongoDBLayer) CountMovieRecords() int {
	cli := mgoLayer.client

	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())

	sess, err := cli.StartSession(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer sess.EndSession(context.TODO())

	var count int64
	// Call WithSession to start a transaction within the new session.
	err = mongo.WithSession(
		context.TODO(),
		sess,
		func(sessCtx mongo.SessionContext) error {
			// Use sessCtx as the Context parameter for InsertOne and FindOne so
			// both operations are run under the new Session.

			if err := sess.StartTransaction(); err != nil {
				return err
			}

			filter := bson.M{}

			moviesCollection := cli.Database(DATABASE).Collection(MOVIES)
			count, err = moviesCollection.CountDocuments(sessCtx, filter)
			// defer cur.Close(ctx)

			if err != nil {
				log.Println(err)
				return err
			}
			return sess.CommitTransaction(context.Background())
		})

	// fmt.Println("ID0:", results[0].Username)
	return int(count)
}

func (mgoLayer *MongoDBLayer) RemoveMovieByID(id string) error {
	cli := mgoLayer.client

	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())

	sess, err := cli.StartSession(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer sess.EndSession(context.TODO())
	// var count []byte

	// Call WithSession to start a transaction within the new session.
	err = mongo.WithSession(
		context.TODO(),
		sess,
		func(sessCtx mongo.SessionContext) error {
			// Use sessCtx as the Context parameter for InsertOne and FindOne so
			// both operations are run under the new Session.

			if err := sess.StartTransaction(); err != nil {
				return err
			}

			newId := bson.M{"_id": id}
			// newId := bson.ObjectIdHex(id)
			moviesCollection := cli.Database(DATABASE).Collection(MOVIES)
			_, err := moviesCollection.DeleteOne(sessCtx, newId)
			if err != nil {
				log.Println(err)
				return err
			}

			// count, _ = json.Marshal(res.DeletedCount)
			return sess.CommitTransaction(context.Background())
		})

	return nil
}

func getNewClient(connection string) (*mongo.Client, error) {
	return mongo.NewClient(options.Client().ApplyURI(connection))
}

func getConnectionURI(db string, user string, pwd string, cluster string) (string, error) {
	if db == "" || user == "" || pwd == "" || cluster == "" {
		return "", fmt.Errorf("Missing Parameters. Please retry!")
	}

	fmt.Println("DB: ", db)
	fmt.Println("User: ", user)
	fmt.Println("Password: ", pwd)
	fmt.Println("Cluster: ", cluster)

	return fmt.Sprintf(`mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority`, user, pwd, cluster, db), nil
}
