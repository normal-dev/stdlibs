import mongoDB from 'mongodb'

let mongoDbUri = process.env.MONGO_DB_URI
if (!mongoDbUri) {
  mongoDbUri = 'mongodb://localhost:27017'
}
const mongoClient = await mongoDB.MongoClient
  .connect(mongoDbUri, {})
  .catch(error => {
    console.error(error)
    process.exit(1)
  })

export default mongoClient
