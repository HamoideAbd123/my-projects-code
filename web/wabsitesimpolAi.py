from bs4 import BeautifulSoup
import requests
import os

def generate_pdf():
    # Get the content of the HTML
    html =
requests.get('http://localhost:8000/').textrequests.get('http://localhost:800/').text

    # Parse the HTML and find all anchor
tags
    soup = BeautifulSoup(html,
'html.parser')
    anchors = soup.find_all('a')

    # Generate the PDF pages from the
anchor tags
    for tag in anchors:
        href = tag['href']
        filepath =
os.path.join(os.getcwd(), 'outcomes/' +
href)

        try:
            with open(filepath, 'wb') as
f:
                response =
requests.get(href, stream=True)
                if response.status_code ==
200:

f.write(response.content)
        except Exception as e:
            print('Error retrieving:',
filepath)
            print(e)

generate_pdf()
