import mongoDB from 'mongodb'

let mongoDbUri = process.env.MONGO_DB_URI
if (!mongoDbUri) {
  console.debug(`can't find MongoDB URI, falling back to %s`, 'mongodb://localhost:27017')
  mongoDbUri = 'mongodb://localhost:27017'
}
const mongoClient = await mongoDB.MongoClient
  .connect(mongoDbUri, {})
  .catch(error => {
    console.error(error)
    process.exit(1)
  })

export default mongoClient
