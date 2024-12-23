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

# Copy custom Nginx configuration files
COPY default.conf /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/nginx.conf

# Create directories and set permissions
RUN mkdir -p /var/cache/nginx/client_temp /var/run/nginx && \
    chmod -R 755 /var/cache/nginx /var/run/nginx

# Ensure proper permissions for non-root user
RUN chmod -R 755 /usr/share/nginx/html && \
    chmod -R 755 /etc/nginx  

# Expose port 3000 for the frontend
EXPOSE 3000

# Run Nginx server
CMD ["nginx", "-g", "daemon off;"]
