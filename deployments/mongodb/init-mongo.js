// MongoDB initialization script
db = db.getSiblingDB('microservices_db');

// Create collections
db.createCollection('users');
db.createCollection('products');

// Create indexes for users collection
db.users.createIndex({ "email": 1 }, { unique: true });
db.users.createIndex({ "created_at": 1 });
db.users.createIndex({ "name": "text", "email": "text" });

// Create indexes for products collection
db.products.createIndex({ "sku": 1 }, { unique: true });
db.products.createIndex({ "category": 1 });
db.products.createIndex({ "created_at": 1 });
db.products.createIndex({ "name": "text", "description": "text" });

// Insert sample data
db.users.insertMany([
    {
        name: "John Doe",
        email: "john@example.com",
        phone: "+1234567890",
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        name: "Jane Smith",
        email: "jane@example.com",
        phone: "+1234567891",
        created_at: new Date(),
        updated_at: new Date()
    }
]);

db.products.insertMany([
    {
        name: "Laptop",
        description: "High-performance laptop for developers",
        price: 1299.99,
        quantity: 50,
        category: "Electronics",
        sku: "LAP001",
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        name: "Mouse",
        description: "Wireless optical mouse",
        price: 29.99,
        quantity: 200,
        category: "Electronics",
        sku: "MOU001",
        created_at: new Date(),
        updated_at: new Date()
    }
]);

print('Database initialization completed successfully!');