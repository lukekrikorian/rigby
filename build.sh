CSS=$(cat /dev/urandom | tr -cd 'a-f0-9' | head -c 16) # Random hash

css () {
	rm ./static/dist/*.css
	cp ./client/index.css ./static/dist/$CSS.css
	cat ./client/css/*.css >> ./static/dist/$CSS.css
}

js () {
	rm ./static/dist/*.js
	npm run $1
}

prod () {
	css
	js "build-prod"
}

printf "[BUILD] "

if [ !$1 ]; then
	if [[ $1 == "--css" ]]; then
		echo "Building CSS only"
		css
		exit 0
	elif [[ $1 == "--js" ]]; then
		echo "Building JS only"
		js "build"
		exit 0
	elif [[ $1 == "--prod" ]]; then
		echo "Building for production"
		prod
		exit 0
	fi 
fi

echo "Building CSS & JS"
css 
js "build"
exit 0