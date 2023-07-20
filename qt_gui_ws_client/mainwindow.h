#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include <QWebSocket>

QT_BEGIN_NAMESPACE
namespace Ui { class MainWindow; }
QT_END_NAMESPACE

class MainWindow : public QMainWindow
{
    Q_OBJECT

public:
    MainWindow(QWidget *parent = nullptr);
    ~MainWindow();

private slots:
    void on_pushBtnCon_clicked();

    void on_pushBtnSend_clicked();

private:
    Ui::MainWindow *ui;
    QWebSocket* ws;
};
#endif // MAINWINDOW_H
