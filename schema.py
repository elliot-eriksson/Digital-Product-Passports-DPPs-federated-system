import  graphene
import graphene
from graphene.relay import Node
from graphene_mongo import MongoengineConnectionField, MongoengineObjectType

from models import comapanies as companiesModel

class companies(MongoengineObjectType):
    class Meta:
        name: "SSAB"
        model: companiesModel