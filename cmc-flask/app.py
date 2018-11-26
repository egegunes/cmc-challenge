import os

from flask import Flask, render_template

from models import db, Currency


POSTGRESQL = {
    "HOST": os.getenv("DBHOST"),
    "NAME": os.getenv("DBNAME"),
    "USER": os.getenv("DBUSER"),
    "PASSWORD": os.getenv("DBPASSWORD")
}

app = Flask(__name__)
app.config['SQLALCHEMY_DATABASE_URI'] = 'postgresql://{USER}:{PASSWORD}@{HOST}/{NAME}'.format(**POSTGRESQL)
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False
db.init_app(app)


@app.route("/")
def get_prices():
    currencies = Currency.query.distinct(Currency.symbol)\
                               .order_by(Currency.symbol, Currency.time.desc())

    return render_template("prices.html", currencies=currencies)
