from fastapi import FastAPI
app = FastAPI()


@app.get("/add")
def add(num1: int, num2: int):
    return num1 + num2
