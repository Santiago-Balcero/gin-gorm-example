# Use the official MySQL image from the Docker Hub
FROM mysql:8.0

# Set environment variables
ENV MYSQL_DATABASE=dogs
ENV MYSQL_ROOT_PASSWORD=12345

# Expose port 3306 to access MySQL
EXPOSE 3306