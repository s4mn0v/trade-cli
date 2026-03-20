# Production

When your SDK is stable and pushed to GitHub:

    Remove the replace line from go.mod.

    Run go get github.com/s4mn0v/bitget@main.

    The CLI will then download the SDK directly from your GitHub repository like any other standard Go library.
