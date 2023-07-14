"""
Logger of proxy queries

"""

# pylint: disable=consider-using-with,too-few-public-methods

import datetime

class Logger:

    """
    Generic logger.
    For specific loggers, _shorten_query() should be rewritten.
    """

    def __init__(self, filename_access, filename_errors):

        self._filename_access = filename_access
        self._filename_errors = filename_errors
        self._log_access = open(filename_access, "a", encoding="utf-8")
        self._log_errors = open(filename_errors, "a", encoding="utf-8")

    def _shorten_query(self, query):
        return query

    def log(self, query, error):
        """
        Log `query` and `error`
        """

        message = str(datetime.datetime.now())
        query = self._shorten_query(query)
        if error != "":
            message += " ERR " + query + " " + error
            self._log_errors.write(message+"\n")
        else:
            message +=  " OK  " + query
            self._log_access.write(message+"\n")


class LoggerWWO(Logger):
    """
    WWO logger.
    """

    def _shorten_query(self, query):
        return "".join([x for x in query.split("&") if x.startswith("q=")])
