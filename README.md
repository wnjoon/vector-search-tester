# vector-search-tester

This is a simple vector search tester. 

## Server

Sentence-bert is executed as a local server. This server is written in python.

To execute server, some binaries are required. 

```bash
$ pip3 install -U Flask sentence-transformers
```

And then, execute server with
```bash
$ python3 server/embed_server.py
```

Server is running at port 6600. You can change port number by modifying `server/embed_server.py`

```python
if __name__ == '__main__':
    # You can change port number in here
    app.run(host='0.0.0.0', port=6600)
```

Sentence-Bert uses `all-MiniLM-L6-v2` model by default. This model is for English text. 

If you want to use Korean model, you can change model to `jhgan/ko-sroberta-multask`. 

```python
# model = SentenceTransformer('jhgan/ko-sroberta-multask') 
model = SentenceTransformer('all-MiniLM-L6-v2')
```

## Client

Client is written in golang. 

### Sentence-Bert

You can run sentence-bert to embed text without additioanl configuration. 

Only you need to do is run sentence-bert server using `python3 server/embed_server.py`. 

Testcode is written in `client/pkg/embedding/sentence_bert_test.go`

> Sentence-Bert creates 384 dimensions vector for English text and 768 dimensions vector for Korean text.

### Gemini

To embed text using gemini, you need to set up environment variable in `.env` file.

`.env` file should be located in root directory of go project (`client` directory). You should add `GEMINI_API_KEY` to `.env` file and set it to your gemini api key.

Testcode is written in `client/pkg/embedding/gemini_test.go`

> Gemini creates 3072 dimensions vector for both English text and Korean text since it is designed to represent 3048 dimensions vector for any language.
