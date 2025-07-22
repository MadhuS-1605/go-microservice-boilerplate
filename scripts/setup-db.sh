#!/bin/bash

set -e

echo "Setting up MongoDB collections and indexes..."

# MongoDB setup commands
mongosh --host localhost:27017 --eval "
use microservices_db;

// Create users collection with indexes
db.users.createIndex({ 'email': 1 }, { unique: true });
db.users.createIndex({ 'created_at': 1 });
db.users.createIndex({ 'name': 'text', 'email': 'text' });

// Create products collection with indexes
db.products.createIndex({ 'sku': 1 }, { unique: true });
db.products.createIndex({ 'category': 1 });
db.products.createIndex({ 'created_at': 1 });
db.products.createIndex({ 'name': 'text', 'description': 'text' });

print('Database setup complete!');
"

echo "Database setup complete!"