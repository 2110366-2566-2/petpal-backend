# generate doc with swagger
# Usage: ./generate_doc.sh

echo "Generating doc with swagger"

# generate doc
cd src/ &&
swag init --parseInternal --parseDependency  -g main.go

echo "Doc generated successfully with pp (again!)"