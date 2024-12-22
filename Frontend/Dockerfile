# Use official Node.js image as base
FROM node:18 AS build

# Set working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm cache clean --force
RUN npm install

# Copy the rest of the application
COPY . .

# Build the Angular app in production mode
RUN npm run build --prod

# Use Nginx to serve the built Angular app
FROM nginx:alpine

# Copy the build output from the previous stage
COPY --from=build /app/dist/my-angular-app/* /usr/share/nginx/html

# Expose port 80 for the frontend
EXPOSE 80

# Run Nginx server
CMD ["nginx", "-g", "daemon off;"]
