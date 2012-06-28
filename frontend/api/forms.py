# -*- coding: utf-8 -*-
# forms.py
# date: 2012-06-29
# author: notedit


import re
from django import forms


class EmailField(forms.CharField):

	def __init__(self):
		forms.CharField.__init__(self,max_length=40)
		return

	def clean(self,value):
		value = forms.CharField.clean(self,value)
		value = value.lower()
		return value


class SignupForm(forms.Form):
	email = EmailField()
	password = forms.CharField()
	nickname = forms.CharField()


class LoginForm(forms.Form):
	email = EmailField()
	password = forms.CharField(max_length=40)
	remember = forms.BooleanField(required=False)

class VerifyNicknameForm(forms.Form):
	nickname = forms.CharField(max_length=40)

class VerifyEmailForm(forms.Form):
	email = EmailField()


	