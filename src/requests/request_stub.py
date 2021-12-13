"""Makes and handles requests to and from the server specified."""
import click


class Request_Stub(object):
    def __init__(self) -> None:
        pass

    """Fake 'request' to the 'server', just returns a string value"""
    def response(self, url, data: str):
        if (data=='BigApple'):
            return ['Markle', 'Apple', 'Markup']
        elif (data=='StripeScore'):
            return ['StripeScoreCore', 'Apple', 'Markup', 'Picket']
        else:
            click.echo(f'Project "{data}" not found!')
            return []

    """Makes a 'request' to a 'server', for now just returns a string value"""
    def make_request(self, url, data):
        # Part where the actual server would be communicated with
        return self.response(url, data)