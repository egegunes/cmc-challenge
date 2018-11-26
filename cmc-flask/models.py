from flask_sqlalchemy import SQLAlchemy


db = SQLAlchemy()


class Currency(db.Model):
    __tablename__ = 'currencies'

    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String)
    symbol = db.Column(db.String)
    price = db.Column(db.Float)
    time = db.Column(db.DateTime)

    def serialize(self):
        return {
            "id": self.id,
            "name": self.name,
            "symbol": self.symbol,
            "price": self.price,
            "time": self.time.strftime("%Y-%m-%dT%H:%M:%S")
        }
