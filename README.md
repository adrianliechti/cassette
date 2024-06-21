# Cassette

Cassette is a simple to host Web Session Recorder and Player.

### Run

```bash
docker run -it --rm -p 3000:3000 ghcr.io/adrianliechti/cassette
```

Open [http://localhost:3000](http://localhost:3000)


### Integration 

```html
<html>
  <head>
    <script src="http://localhost:3000/cassette.min.cjs"></script>
  </head>
<html>
```