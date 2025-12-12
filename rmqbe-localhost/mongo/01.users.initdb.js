db = db.getSiblingDB("rmqbe");

// Plaintext passwords (for reference):
// mqtt1 -> password1
// mqtt2 -> password2

db.createCollection("users");

db.users.insertMany([
  {
    username: "mqtt1",
    password: "$2a$10$v1Gd1slX0wzN0t7Um1rzSOWx8WtIdyQn2PvA4dPaqw/ENuUz4xGdm",
    vhosts: ["/"],
    permissions: [
      { resource: "exchange", name: "amq.topic", permission: "read" },
      { resource: "exchange", name: "amq.topic", permission: "write" },
      { resource: "queue", name: "#", permission: "read" },
      { resource: "queue", name: "#", permission: "write" },
    ],
  },
  {
    username: "mqtt2",
    password: "$2a$10$4WrQxJjZ0CqPf0Em2zhqPeHh0jK7PWJ5mnw9nwzN2Z6F5h.3P4gdi",
    vhosts: ["/"],
    permissions: [
      { resource: "exchange", name: "amq.topic", permission: "read" },
      { resource: "exchange", name: "amq.topic", permission: "write" },
      { resource: "queue", name: "#", permission: "read" },
      { resource: "queue", name: "#", permission: "write" },
    ],
  },
]);
